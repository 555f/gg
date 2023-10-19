package httperror

// func extractErrorValues(methods []*types.Func) (statusCode int64, errorText, code string, dataType any, ok bool) {
// 	var (
// 		statusCodeMethod, codeMethod, errorMethod, dataMethod *types.Func
// 	)
// 	for _, method := range methods {
// 		switch method.Name {
// 		case "Data":
// 			dataMethod = method
// 		case "StatusCode":
// 			statusCodeMethod = method
// 		case "Error":
// 			errorMethod = method
// 		case "Code":
// 			codeMethod = method
// 		}
// 	}
// 	if statusCodeMethod == nil || errorMethod == nil {
// 		return 0, "", "", nil, false
// 	}
// 	statusCode = statusCodeMethod.Returns[0].(int64)
// 	if dataMethod != nil {
// 		if len(dataMethod.Returns) == 1 {
// 			dataType = dataMethod.Returns[0]
// 		}
// 	}
// 	if errorMethod != nil {
// 		errorText, _ = errorMethod.Returns[0].(string)
// 	}
// 	if codeMethod != nil {
// 		code, _ = codeMethod.Returns[0].(string)
// 	}
// 	ok = true
// 	return
// }
