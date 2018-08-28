package handler

//模板: 普通函数
const TEMP_FUNC = `
            //{{.funcnote}}
            func {{.funcname}}(req {{.funcname}}Req) ({{.funcname}}Resp, error) {
              resp := {{.funcname}}Resp{}
              {{- if .body}}
                {{.body}}
              {{- end}}
              return resp, nil
            }`

//模板: handler的主函数
const TEMP_HANDLER = `package handler

            //##################################################
            //                   主处理函数
            //##################################################
            {{- range .mainFuncCode}}
            {{.}}
            {{- end}}

            //##################################################
            //                     过程函数
            //##################################################
            {{- range .subFuncCodeList}}
            {{.}}
            {{- end}}

            //##################################################
            //                  封装中间请求数据
            //##################################################
            {{- range .innerDataDeclareCode -}}
            {{.}}
            {{- end}}
            {{- range .reqMakerCodeList}}
            {{.}}
            {{- end}}
            {{- range .makeResponse}}
            {{.}}
            {{- end}}`

//模板: handler的body
const TEMP_HANDLER_BODY = `//{{.funcnote}}
            innerData.req{{.funcname}} = innerData.make{{.funcname}}Req()
            resp{{.funcname}}, err := {{.funcname}}(innerData.req{{.funcname}})
            if err != nil {
              return {{.mainFunc}}Resp{}, err
            }
            innerData.resp{{.funcname}} = resp{{.funcname}}`

//模板: handler的数据流集
const TEMP_HANDLER_INNER_DATA_INIT = `//单个请求涉及的中间数据集合
            type innerData{{- range .mainFunc}}{{.}}{{- end}} struct {
              {{- range .mainFunc}}
              req {{.}}Req
              // resp {{.}}Resp   //(no need)
              {{- end}}
              {{- range .subFuncList}}
              req{{.}} {{.}}Req
              resp{{.}} {{.}}Resp
              {{- end}}
            }`

//模板: handler的数据流集初始赋值
const TEMP_HANDLER_INNER_DATA_DEFINE = `innerData := innerData{{.mainFunc}}{
                req:  req,
            }`

//模板: handler的结果组装
const TEMP_HANDLER_MAKE_RESP = `//组装返回数据
            func (*innerData{{.mainFunc}}) makeResp() {{.mainFunc}}Resp {
              return {{.mainFunc}}Resp{}
            }`

//模板: proxy的body
const TEMP_PROXY = `//request
            {{- if .methodget}}
            url := fmt.Sprintf("http://%s{{.reqpath}}?%s", ReqAddr, proxy.Struct2Querystr(req))
            httpcode, body := proxy.Get(url)
            {{- end}}
            {{- if .methodpost}}
            url := fmt.Sprintf("http://%s{{.reqpath}}", ReqAddr)
            httpcode, body := proxy.PostJson(url, req)
            {{- end}}
            if httpcode != http.StatusOK {
              return nil, errors.New("request failed.")
            }

            //response
            err := json.Unmarshal(body, rst)
            if err != nil {
              return nil, err
            }`

//模板: 数据流涉及的请求封装
const TEMP_REQ_MAKER = `//组装{{.funcname}}的请求数据
            func (innerData *innerData{{.mainFunc}}) make{{.funcname}}Req() ({{.funcname}}Req) {
              return {{.funcname}}Req{}
            }`

//模板: mock函数
const TEMP_MOCK_FUNC = `//mock {{.funcname}}
            func mock{{.funcname}}(req {{.funcname}}Req) *{{.funcname}}Resp {
              rst := &{{.funcname}}Resp{}

              {{.body}}
              return rst
            }
            `
