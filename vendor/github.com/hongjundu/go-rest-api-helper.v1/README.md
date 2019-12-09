# golang rest api helper

A set of small data structure that help to write rest API response.

* OK response format

        {
            "status": "ok",
            "data":
            {
                "firstName": "Jon",
                "lastName": "Snow"
            }
        }

* Error response format

        {
            "status": "error",
            "code": some_error_code,
            "msg": "some error message"
        }

## Example

        package main

        import (
            "encoding/json"
            "github.com/hongjundu/go-rest-api-helper.v1"
            "log"
            "net/http"
        )

        func main() {
            log.Printf("server starts...")

            http.HandleFunc("/ok", wrapHandler(handlerReturnsOK))
            http.HandleFunc("/error", wrapHandler(handlerReturnsError))

            log.Fatal(http.ListenAndServe(":8080", nil))

            log.Printf("server exists")
        }

        func handlerReturnsOK(r *http.Request) (response interface{}, err error) {
            type Resp struct {
                FirstName string `json:"firstName"`
                LastName  string `json:"lastName"`
            }
            response = apihelper.NewOKResponse(Resp{FirstName: "Jon", LastName: "Snow"})
            return
        }

        func handlerReturnsError(r *http.Request) (response interface{}, err error) {
            err = apihelper.NewError(http.StatusBadRequest, "bad request")
            return
        }

        func wrapHandler(handler func(r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
            return func(w http.ResponseWriter, r *http.Request) {
                var err error
                var response interface{}
                var responseBody []byte

                if response, err = handler(r); err == nil {
                    responseBody, err = json.Marshal(response)
                }

                if err != nil {
                    if apiErr, ok := err.(apihelper.ApiError); ok {
                        w.WriteHeader(apiErr.Code())
                    }

                    responseBody, _ = json.Marshal(apihelper.NewErrorResponse(err))
                }

                w.Write(responseBody)
            }
        }


