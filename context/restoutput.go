package context

import (
    "compress/flate"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
)

type RESTError struct {
    Massage string            `json:"massage"`
    Errors  map[string]string `json:"errors"`
}

// render default application error page with error and stack string.
func (output *BeegoOutput) RESTErr(err interface{}, r *http.Request) {
    errors := make(map[string]string)
    errors["method"] = r.Method
    errors["resource"] = r.RequestURI
    errors["code"] = "inner_error"

    re := RESTError{
        Massage: fmt.Sprint(err),
        Errors:  errors,
    }
    output.RESTJson(http.StatusInternalServerError, re, true, true)
}

func (output *BeegoOutput) RESTJson(status int, data interface{}, hasIndent bool, coding bool) error {
    output.Header("Content-Type", "application/json;charset=UTF-8")
    var content []byte
    var err error
    if hasIndent {
        content, err = json.MarshalIndent(data, "", "  ")
    } else {
        content, err = json.Marshal(data)
    }
    if err != nil {
        http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
        return err
    }
    if coding {
        content = []byte(stringsToJson(string(content)))
    }
    output_writer := output.Context.ResponseWriter.(io.Writer)
    if output.EnableGzip == true && output.Context.Input.Header("Accept-Encoding") != "" {
        splitted := strings.SplitN(output.Context.Input.Header("Accept-Encoding"), ",", -1)
        encodings := make([]string, len(splitted))

        for i, val := range splitted {
            encodings[i] = strings.TrimSpace(val)
        }
        for _, val := range encodings {
            if val == "gzip" {
                output.Header("Content-Encoding", "gzip")
                output_writer, _ = gzip.NewWriterLevel(output.Context.ResponseWriter, gzip.BestSpeed)

                break
            } else if val == "deflate" {
                output.Header("Content-Encoding", "deflate")
                output_writer, _ = flate.NewWriter(output.Context.ResponseWriter, flate.BestSpeed)
                break
            }
        }
    } else {
        output.Header("Content-Length", strconv.Itoa(len(content)))
    }
    output.Context.ResponseWriter.WriteHeader(status)
    output_writer.Write(content)
    switch output_writer.(type) {
    case *gzip.Writer:
        output_writer.(*gzip.Writer).Close()
    case *flate.Writer:
        output_writer.(*flate.Writer).Close()
    }
    return nil
}