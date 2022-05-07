package response

import (
	"blockchain/pkg/customerrors"
	"blockchain/pkg/logger"
	"blockchain/pkg/network"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SendCsv(ctx context.Context, w http.ResponseWriter, fileName string, bom bool, records [][]string) {
	// text/csv だと ブラウザにTextで表示されていまうので OctetStreamに固定
	w.Header().Set(network.ContentType, network.ApplicationOctetStream)
	w.Header().Set(network.ContentDisposition, fmt.Sprintf("attachment;filename=\"%s\";filename*=UTF-8''%s", fileName, url.QueryEscape(fileName)))
	if bom {
		if _, err := w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
			logger.Logging.Error("failed to write ", zap.Error(err))
			return
		}
	}
	cw := csv.NewWriter(w)
	w.WriteHeader(http.StatusOK)
	var err error
	if err = cw.WriteAll(records); err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			logger.Logging.Info(fmt.Sprintf("connection reset by peer. %s", getErrorFileMessage(network.ApplicationOctetStream, &fileName, nil)))
			return
		}
		if strings.Contains(err.Error(), "broken pipe") {
			logger.Logging.Info(fmt.Sprintf("broken pipe. %s", getErrorFileMessage(network.ApplicationOctetStream, &fileName, nil)))
			return
		}
		logger.Logging.Error("failed to write", zap.Error(err))
	}
}

func Send(ctx context.Context, status int, msg interface{}, w http.ResponseWriter, opts ...SendOption) {
	dopts := &SendOptions{}
	for i := range opts {
		opts[i](dopts)
	}
	logOrgErr(ctx, dopts.OrgErr)

	w.Header().Set(network.ContentType, network.ApplicationJson)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		logger.Logging.Error("failed to write", zap.Error(err))
	}
}

func logOrgErr(ctx context.Context, orgErr error) {
	if orgErr == nil {
		return
	}

	var validationErrs validator.ValidationErrors
	var customErr *customerrors.CustomError
	switch {
	case errors.As(orgErr, &validationErrs):
		logger.Logging.Warn("validation error", zap.Error(validationErrs))
	case errors.As(orgErr, &customErr):
		if customerrors.IsServerError(customErr.Code()) {
			logger.Logging.Error("server custom error", zap.Object("error", customErr))
		} else {
			logger.Logging.Warn("client custom error", zap.Object("error", customErr))
		}
	default:
		logger.Logging.Error("error default case", zap.Error(orgErr))
	}
}

func SendFile(ctx context.Context, w http.ResponseWriter, closer io.ReadCloser, contentType string, fileName, responseRange *string) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, closer); err != nil {
		InternalServerError(ctx, w, WithError(err))
		return
	}
	defer buf.Reset()

	w.Header().Set(network.ContentType, contentType)
	if fileName != nil {
		w.Header().Set(network.ContentDisposition, fmt.Sprintf("attachment;filename=\"%s\";filename*=UTF-8''%s", *fileName, url.QueryEscape(*fileName)))
	}
	if responseRange != nil && *responseRange != "" {
		w.Header().Set(network.ContentRange, *responseRange)
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		if strings.Contains(err.Error(), "connection reset by peer") {
			logger.Logging.Info(fmt.Sprintf("connection reset by peer. %s", getErrorFileMessage(contentType, fileName, responseRange)))
			return
		}
		if strings.Contains(err.Error(), "broken pipe") {
			logger.Logging.Info(fmt.Sprintf("broken pipe. %s", getErrorFileMessage(contentType, fileName, responseRange)))
			return
		}
		logger.Logging.Error("failed to write", zap.Error(err))
	}
}

func getErrorFileMessage(contentType string, fileName, responseRange *string) string {
	values := "contentType: " + contentType
	if fileName != nil {
		values += ",fileName: " + *fileName
	}
	if responseRange != nil {
		values += ",responseRange: " + *responseRange
	}
	return values
}

func Redirect(ctx context.Context, w http.ResponseWriter, location string, opts ...SendOption) {
	dopts := &SendOptions{}
	for i := range opts {
		opts[i](dopts)
	}
	logOrgErr(ctx, dopts.OrgErr)

	w.Header().Set(network.Location, location)
	w.WriteHeader(http.StatusFound)
}

func RedirectWithHTML(ctx context.Context, w http.ResponseWriter, location string) {
	OKWithHTML(
		ctx,
		w,
		fmt.Sprintf(`<!DOCTYPE html><html><head><meta http-equiv="refresh" content="0;URL=%s"></head></html>`, html.EscapeString(location)),
	)
}

func OKWithHTML(ctx context.Context, w http.ResponseWriter, htmlStr string, opts ...SendOption) {
	dopts := &SendOptions{}
	for i := range opts {
		opts[i](dopts)
	}
	logOrgErr(ctx, dopts.OrgErr)

	w.Header().Set(network.ContentType, network.TextHtml)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(htmlStr)); err != nil {
		logger.Logging.Error("failed to write", zap.Error(err))
	}
}

func OKWithMsg(ctx context.Context, w http.ResponseWriter, msg interface{}, opts ...SendOption) {
	Send(ctx, http.StatusOK, msg, w, opts...)
}

func OK(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusOK, getMessage(http.StatusOK), w, opts...)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusInternalServerError, getMessage(http.StatusInternalServerError), w, opts...)
}

func Error(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusInternalServerError, getMessage(http.StatusInternalServerError), w, opts...)
}

func BadRequest(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusBadRequest, getMessage(http.StatusBadRequest), w, opts...)
}

func Forbidden(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusForbidden, getMessage(http.StatusForbidden), w, opts...)
}

func NotFound(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusNotFound, getMessage(http.StatusNotFound), w, opts...)
}

func Unauthorized(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusUnauthorized, getMessage(http.StatusUnauthorized), w, opts...)
}

func Conflict(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusConflict, getMessage(http.StatusConflict), w, opts...)
}

func PayloadTooLarge(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusRequestEntityTooLarge, getMessage(http.StatusRequestEntityTooLarge), w, opts...)
}

func RangeNotSatisfiable(ctx context.Context, w http.ResponseWriter, opts ...SendOption) {
	Send(ctx, http.StatusRequestedRangeNotSatisfiable, getMessage(http.StatusRequestedRangeNotSatisfiable), w, opts...)
}

func ErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	var customErr *customerrors.CustomError
	if errors.As(err, &customErr) {
		switch customErr.Code() {
		case customerrors.ErrCodeConflict,
			customerrors.ErrCodeOptimisticLock:
			Conflict(ctx, w, WithError(err))
		case customerrors.ErrCodeInvalidArgument:
			BadRequest(ctx, w, WithError(err))
		case customerrors.ErrCodeNotFound:
			NotFound(ctx, w, WithError(err))
		case customerrors.ErrUnAuthorized:
			Unauthorized(ctx, w, WithError(err))
		case customerrors.ErrForbidden:
			Forbidden(ctx, w, WithError(err))
		case customerrors.ErrPayloadTooLarge:
			PayloadTooLarge(ctx, w, WithError(err))
		case customerrors.ErrRangeNotSatisfiable:
			RangeNotSatisfiable(ctx, w, WithError(err))
		case customerrors.ErrUnexpected:
			InternalServerError(ctx, w, WithError(err))
		default:
			InternalServerError(ctx, w, WithError(err))
		}
		return
	}
	logger.Logging.Error("err isn't CustomError", zap.Error(err))
	InternalServerError(ctx, w, WithError(err))
}

type errorMessage struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func getMessage(code int) errorMessage {
	return errorMessage{
		Status:  code,
		Message: http.StatusText(code),
	}
}
