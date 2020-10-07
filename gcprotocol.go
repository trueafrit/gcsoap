// gcprotocol.go
package gcsoap

import (
	"bytes"
	"encoding/xml"
	"errors"

	"fmt"
	//"gouploader/Runconfig"
	"strconv"

	logrus "github.com/sirupsen/logrus"
)

var Runconfig Config

type Config interface {
	GetString(string) string
	SetString(string, string)
	GetInt(string) int
	SetInt(string, int)
}

type GCGetToken struct {
	XMLName     xml.Name `xml:"cbd2:GetToken"`
	FormerToken string   `xml:"cbd2:formerToken"`
}

type GCGetDataTable struct {
	XMLName   xml.Name `xml:"cbd2:GetDataTable"`
	Token     string   `xml:"cbd2:token"`
	DataTable string   `xml:"cbd2:dataTable"`
}

/*GetAgreedSegments() Веб-сервис: AgreedRoutes.asmx */
type GCGetAgreedSegments struct {
	XMLName xml.Name `xml:"cbd2:GetAgreedSegments"`
	Token   string   `xml:"cbd2:token"`
	Options string   `xml:"cbd2:options"`
}

type GCAgreedSegment struct {
	XMLName       xml.Name `xml:"cbd2:SendAgreedSegment"`
	Id            int      `xml:"cbd2:id" json:"id"`                       // Идентификатор изменяемой записи или -1 в случае, когда требуется создать новую запись
	AirwayCode    string   `xml:"cbd2:airwayCode" json:"airwayCode"`       // Код ИКАО трассы (лат.)
	AtsBeginPoint string   `xml:"cbd2:atsBeginPoint" json:"atsBeginPoint"` // Код ИКАО начальной точки (лат.)
	AtsEndPoint   string   `xml:"cbd2:atsEndPoint" json:"atsEndPoint"`     // Код ИКАО крайней точки (лат.)
	IsOneWay      int      `xml:"cbd2:isOneWay" json:"isOneWay"`           // Признак однонаправленности: 1 – участок однонаправленный в прямом направлении; 2 – участок однонаправленный в обратном направлении
	LevelFrom     string   `xml:"cbd2:levelFrom" json:"levelFrom"`         // Эшелон от
	LevelTo       string   `xml:"cbd2:levelTo" json:"levelTo"`             // Эшелон до
	DateFrom      string   `xml:"cbd2:dateFrom" json:"dateFrom"`           // Дата и время начала использования в формате “YYYY MM DD"T"HH24:MI"Z"”
	DateTo        string   `xml:"cbd2:dateTo" json:"dateTo"`               // Дата и время окончания использования в формате “YYYY MM DD"T"HH24:MI"Z"”
	ExceptionMsg  string   `xml:"cbd2:exceptionMsg" json:"exceptionMsg"`   // Примечание (максимальная длина – 2000 символов)
}

// --- for parse XML output from SOAP function GetAgreesSegmets
type GCEnvelop struct {
	XMLName xml.Name `xml:"Envelope" json:"-"`
	Body    GCBoby   `xml:"Body" json:"-"`
}

type GCBoby struct {
	XMLName  xml.Name                    `xml:"Body" json:"-"`
	Response GCGetAgreedSegmentsResponse `xml:"GetAgreedSegmentsResponse" json:"-"`
}

type GCGetAgreedSegmentsResponse struct {
	XMLName xml.Name                  `xml:"GetAgreedSegmentsResponse" json:"-"`
	Result  GCGetAgreedSegmentsResult `xml:"GetAgreedSegmentsResult" json:"-"`
}

type GCGetAgreedSegmentsResult struct {
	XMLName xml.Name `xml:"GetAgreedSegmentsResult" json:"-"`
	Rowset  GCRowset `xml:"ROWSET" json:"-"`
}

type GCRowset struct {
	XMLName xml.Name `xml:"ROWSET" json:"-"`
	Rows    []GCRow  `xml:"ROW" json:"ROWSET"`
}

/*
<ROW>
<HEAD_ATSREC_ID>652437</HEAD_ATSREC_ID>
<ID>652437</ID>
<ZC_ID>11</ZC_ID>
<ATSRECNO>4</ATSRECNO>
<AIRWCODE>W241</AIRWCODE>
<ATSBEGINPOINT>BD</ATSBEGINPOINT>
<ATSENDPOINT>IKT</ATSENDPOINT>
<ISONEWAY>0</ISONEWAY>
<LEVELFROM>F210</LEVELFROM>
<LEVELTO>F530</LEVELTO>
<ISEXCEPTION>0</ISEXCEPTION>
<RECPROCESS>1</RECPROCESS>
<ATSRECDATE>2020-01-09T00:00:00Z</ATSRECDATE>
<DATEFROM>2020-01-09T05:48:00Z</DATEFROM>
<DATETO>2020-12-31T00:59:00Z</DATETO>
<ISDELETE>1</ISDELETE>
</ROW>
*/
type GCRow struct {
	XMLName        xml.Name `xml:"ROW" json:"-"`
	Head_atsrec_id int      `xml:"HEAD_ATSREC_ID" json:"head_atsrec_id"`
	Id             int      `xml:"ID" json:"id"`       // Идентификатор изменяемой записи или -1 в случае, когда требуется создать новую запись
	Zc_id          int      `xml:"ZC_ID" json:"zc_id"` // Идентификатор зонального центра
	Atsrecno       int      `xml:"ATSRECNO" json:"atsrecno"`
	Airwcode       string   `xml:"AIRWCODE" json:"airwcode"`           // Код ИКАО трассы (лат.)
	AtsBeginPoint  string   `xml:"ATSBEGINPOINT" json:"atsBeginPoint"` // Код ИКАО начальной точки (лат.)
	AtsEndPoint    string   `xml:"ATSENDPOINT" json:"atsEndPoint"`     // Код ИКАО крайней точки (лат.)
	IsOneWay       int      `xml:"ISONEWAY" json:"isOneWay"`           // Признак однонаправленности: 1 – участок однонаправленный в прямом направлении; 2 – участок однонаправленный в обратном направлении
	LevelFrom      string   `xml:"LEVELFROM" json:"levelFrom"`         // Эшелон от
	LevelTo        string   `xml:"LEVELTO" json:"levelTo"`             // Эшелон до
	ISEXCEPTION    int      `xml:"ISEXCEPTION" json:"isexception"`
	RECPROCESS     int      `xml:"RECPROCESS" json:"recprocess"`
	ATSRECDATE     string   `xml:"ATSRECDATE" json:"atsrecdate"`
	DateFrom       string   `xml:"DATEFROM" json:"dateFrom"` // Дата и время начала использования в формате “YYYY MM DD"T"HH24:MI"Z"”
	DateTo         string   `xml:"DATETO" json:"dateTo"`     // Дата и время окончания использования в формате “YYYY MM DD"T"HH24:MI"Z"”
	ISDELETE       int      `xml:"ISDELETE" json:"isdelete"`
}

func (b GCAgreedSegment) String() string {
	return fmt.Sprintf(`GCAgreedSegment ( id: %d AirwayCode: %s AtsBeginPoint: %s AtsEndPoint: %s IsOneWay: %d LevelFrom: %s LevelTo: %s
 DateFrom: %s DateTo: %s ExceptionMsg: %s)`,
		b.Id, b.AirwayCode, b.AtsBeginPoint, b.AtsEndPoint, b.IsOneWay, b.LevelFrom, b.LevelTo, b.DateFrom, b.DateTo, b.ExceptionMsg)
}

type GCDeleteAgreedSegment struct {
	XMLName xml.Name `xml:"cbd2:DeleteAgreedSegment"`
	Id      string   `xml:"cbd2:id"`
}

// формирует запрос
func GC_CreateRequest(login string, password string, bData interface{}) MsoapRequest {
	req := MsoapRequest{
		XMLNsSoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsCBD2:    "http://GKOVD/CBD2",
		Header: MsoapRequestHeader{
			AuthData: MsoapAuthData{
				Login:    login,
				Password: password,
			},
		},
	}

	b := MsoapBody{}
	b.Body = bData
	req.Body = b

	return req
}

/*
Надо обрабатывать не только ошибки вызова http.client.post но и возвращаемый результат

RESPONSE:
 <?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body>
	<SendAgreedSegmentResponse xmlns="http://GKOVD/CBD2">
		<SendAgreedSegmentResult>
			<error code="1003">access denied</error>
		</SendAgreedSegmentResult>
	</SendAgreedSegmentResponse>
</soap:Body>
</soap:Envelope>
*/

// Еще и такие ошибки:
/*
<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body><soap:Fault><faultcode>soap:Server</faultcode>
<faultstring>System.Web.Services.Protocols.SoapException: Серверу не удалось обработать запрос. ---&gt; System.ArgumentNullException: Ссылка на строку не ссылается на экземпляр String.
Имя параметра: s
   в System.DateTimeParse.ParseExact(String s, String format, DateTimeFormatInfo dtfi, DateTimeStyles style)
   в System.DateTime.ParseExact(String s, String format, IFormatProvider provider)
   в GKOVD.CBD2.AgreedRoutes.SendAgreedSegment(Int32 id, String airwayCode, String atsBeginPoint, String atsEndPoint, Int32 isOneWay, String levelFrom, String levelTo, String dateFrom, String dateTo, String exceptionMsg) в D:\prj\VS2010\WebApp\CBD2\1.0\AgreedRoutes.asmx.cs:строка 126
   --- Конец трассировки внутреннего стека исключений ---</faultstring><detail /></soap:Fault></soap:Body></soap:Envelope>

//TODO: Еще и такие
<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body>
<SendAgreedSegmentResponse xmlns="http://GKOVD/CBD2"><SendAgreedSegmentResult>
<result errorCode="-8" xmlns="" /></SendAgreedSegmentResult>
</SendAgreedSegmentResponse>
</soap:Body>
</soap:Envelope>
*/
func checkXMLResponseForError(data []byte) error {
	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "checkXMLResponseForError",
	})

	SendAgreedSegmentErrors := map[int]string{
		-1: `Отсутствуют обязательные параметры. (к ним относятся код трассы, конечная и начальная точка,
		конечный и начальный эшелон, дата начала и окончания действия, признак однонаправленности, номер зонального центра)`,
		-2: `Код трассы задан в неправильном формате. Код трассы может содержать только латинские заглавные буквы и цифры.
		Длина кода трассы не должна превышать 6 символов.`,
		-3: `Начальная точка задана в неправильном формате. Названия начальной и конечной точек могут содержать только латинские
заглавные буквы. Длина названий начальной и конечной точек не должна превышать 6 символов.`,
		-4: `Конечная точка задана в неправильном формате. Названия начальной и конечной точек могут содержать только латинские
заглавные буквы. Длина названий начальной и конечной точек не должна превышать 6 символов.`,
		-5: `Начальный эшелон задан в неправильном формате. Названия эшелонов должны соответствовать формату F999, где на месте
каждой цифры 9 могут находиться любые цифры. Длина названия эшелона должна составлять 4 символа.`,
		-6: `Конечный эшелон задан в неправильном формате. Названия эшелонов должны соответствовать формату F999, где на месте каждой
цифры 9 могут находиться любые цифры. Длина названия эшелона должна составлять 4 символа.`,
		-7: `Конечный эшелон меньше начального. Конечный эшелон должен быть больше начального.`,
		-8: `Нарушение целостности данных. Конечная и начальная дата должны быть больше текущей. Также конечная дата должна быть больше
начальной.`,
		-9: `Попытка изменить запись, закончившую свое действие.`,
		-10: `Попытка изменить запись, начавшую свое действие. В записи, начавшей своей действие, для изменения доступна
только дата окончания действия.`,
		-11: `Начальный эшелон не входит в список допустимых эшелонов.(10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120,
130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 260, 270, 280, 290, 300, 310, 320, 330, 340, 350,
360, 370, 380, 390, 400, 410, 430, 450, 470, 490, 510, 530, 550, 570, 590, 610, 630, 650, 670)`,
		-12:    `Конечный эшелон не входит в список допустимых эшелонов.`,
		-13:    `Номер зонального центра не входит в список допустимых.(1, 4, 5, 7, 8, 11, 13, 17)`,
		-14:    `Признак однонаправленности не входит в список допустимых значений (0,1)`,
		-15:    `Ошибка, что такой сегмент уже есть с такими же данными (предположительно)`,
		-20070: `Ошибка удаления сегмента, сегмент с таким id не существует (предположительно)`,
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))

	var inElement, code string
	var l XMLQuery

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		if se, ok := t.(xml.StartElement); ok {
			inElement = se.Name.Local
			if inElement == "error" || inElement == "faultstring" || inElement == "result" {

				err := decoder.DecodeElement(&l, &se)
				if err != nil {
					log.Error(err)
				}

				if inElement == "faultstring" {
					return errors.New("faultstring : " + l.Loc)
				}

				for _, value := range se.Attr {
					if value.Name.Local == "code" && inElement == "error" {
						code = value.Value
						return errors.New(code + " : " + l.Loc)
					}
					if value.Name.Local == "errorCode" && inElement == "result" {
						code = value.Value
						if res, err := strconv.Atoi(code); err == nil {
							if res == 0 {
								return nil
							}
							return errors.New(code + " : " + SendAgreedSegmentErrors[res])
						}
						return errors.New(code + " : " + "Unknown Error")
					}
				}

			}
		}
	}

	return nil
}

func GC_GetToken(inToken string) (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "GC_GetToken",
	})
	gt := GCGetToken{}
	gt.FormerToken = inToken

	v := GC_CreateRequest(Runconfig.GetString("SoapLogin"), Runconfig.GetString("SoapPassword"), gt)

	request, err := xml.MarshalIndent(v, "", "  ")

	log.Debug(string(request), "\nWITH ERROR: ", err)

	if err != nil {
		return "", err
	}

	response, err := MsoapCall(Runconfig.GetString("MetaUrl"), "http://GKOVD/CBD2/GetToken", request, 30)
	if err != nil {
		return "", err
	}

	// parse xml and get only token

	decoder := xml.NewDecoder(bytes.NewReader(response))

	var inElement, res string
	var l XMLQuery

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		if se, ok := t.(xml.StartElement); ok {
			inElement = se.Name.Local
			if inElement == "GetTokenResult" {

				err := decoder.DecodeElement(&l, &se)
				if err != nil {
					log.Error(err)
				}

				log.Debug("Key: %s value: %s\n", inElement, l.Loc)
				res = l.Loc

			}
		}
	}

	return res, nil

}

// Параметры функции: SendAgreedSegment() Веб-сервис: AgreedRoutes.asmx
//
// Функция предназначена для передачи/изменения в КСА ПВД ГЦ перечня маршрутов ОВД, разрешённых для использования по согласованию.
//
// В случае вызова с идентификатором записи необходимо заполнить только изменяемые параметры новыми значениями; значения неизменяемых параметров следует оставить пустыми.
// В случае успеха функция возвращает идентификатор вставленной или изменённой записи (неотрицательное число).
// В случае возникновения ошибки функция возвращает отрицательное число.

// Нормальный ответ:
/*
<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:xsd="http://www.w3.org/2001/XMLSchema">
<soap:Body>
	<SendAgreedSegmentResponse xmlns="http://GKOVD/CBD2"><SendAgreedSegmentResult>
		<result id="2073865114" recProcess="1" xmlns=""><errorMessage>null</errorMessage>
		</result>
	</SendAgreedSegmentResult>
</SendAgreedSegmentResponse>
</soap:Body>
</soap:Envelope>
*/
func GC_SendAgreedSegment(segm GCAgreedSegment) (int, error) {
	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "GC_SendAgreedSegment",
	})

	log.Debug("SEND TO ", Runconfig.GetString("AgreedRoutesUrl"))

	v := GC_CreateRequest(Runconfig.GetString("SoapLogin"), Runconfig.GetString("SoapPassword"), segm)

	request, err := xml.MarshalIndent(v, "", "  ")

	log.Debug("REQUEST: ", string(request))

	if err != nil {
		log.Error("WITH ERROR: ", err)
		return -1, err
	}

	response, err := MsoapCall(Runconfig.GetString("AgreedRoutesUrl"), "http://GKOVD/CBD2/SendAgreedSegment", request, 60)

	log.Debug("RESPONSE: ", string(request))

	if err != nil {
		log.Debug("WITH ERROR: ", err)
		return -1, err
	}

	err = checkXMLResponseForError(response)

	if err != nil {
		log.Error("RESPONSE:\n", string(response), "\nWITH ERROR", err)
		return -1, err
	}

	// parse xml and get result Id

	decoder := xml.NewDecoder(bytes.NewReader(response))
	var inElement string
	var res int = -99
	var l XMLQuery
	var errRes error = nil

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		if se, ok := t.(xml.StartElement); ok {
			inElement = se.Name.Local
			//SendAgreedSegmentResult
			if inElement == "result" {

				for _, value := range se.Attr {
					if value.Name.Local == "id" {
						log.Debug("Key: %s value: %s\n", inElement, l.Loc)
						if res, err = strconv.Atoi(value.Value); err != nil {
							errRes = errors.New(l.Loc + " не является целым числом ID.")
							log.Debug(l.Loc, " не является целым числом ID.")
						}
					}
				}

			}
			if inElement == "errorMessage" {
				err := decoder.DecodeElement(&l, &se)
				if err != nil {
					log.Error(err)
				}

				if l.Loc != "" || l.Loc != "null" || l.Loc != "Null" || l.Loc != "NULL" {
					errRes = errors.New(l.Loc)
				}
			}
		}
	}

	return res, errRes
}

/*<soapenv:Body>
     <cbd2:DeleteAgreedSegment>
        <cbd2:id>2073743693</cbd2:id>
     </cbd2:DeleteAgreedSegment>
  </soapenv:Body>

result ok:
<soap:Body><DeleteAgreedSegmentResponse xmlns="http://GKOVD/CBD2"><DeleteAgreedSegmentResult>
<result errorCode="0" xmlns="" /></DeleteAgreedSegmentResult></DeleteAgreedSegmentResponse></soap:Body>

result not ok(not exists):
<soap:Body><DeleteAgreedSegmentResponse xmlns="http://GKOVD/CBD2"><DeleteAgreedSegmentResult>
<result errorCode="-20070" xmlns="" /></DeleteAgreedSegmentResult></DeleteAgreedSegmentResponse></soap:Body>

*/
func GC_DeleteAgreedSegment(id int) (bool, error) {
	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "GC_DeleteAgreedSegment",
	})

	ds := GCDeleteAgreedSegment{}
	ds.Id = strconv.Itoa(id)
	v := GC_CreateRequest(Runconfig.GetString("SoapLogin"), Runconfig.GetString("SoapPassword"), ds)

	request, err := xml.MarshalIndent(v, "", "  ")

	log.Debug("REQUEST: ", string(request))

	if err != nil {
		log.Error("WITH ERROR: ", err)

	}

	response, err := MsoapCall(Runconfig.GetString("AgreedRoutesUrl"), "http://GKOVD/CBD2/DeleteAgreedSegment", request, 30)

	log.Debug("RESPONSE: ", string(response))

	if err != nil {
		log.Error("MsoapCall RESPONSE: ", string(response), "WITH ERROR: ", err)

	}

	err = checkXMLResponseForError(response)

	if err != nil {
		return false, err
	}

	return true, nil

}

func GC_GetAgreedSegments(token string) (GCRowset, error) {
	var result GCRowset

	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "GC_GetAgreedSegments",
	})

	ds := GCGetAgreedSegments{}
	ds.Token = token
	ds.Options = ""
	v := GC_CreateRequest(Runconfig.GetString("SoapLogin"), Runconfig.GetString("SoapPassword"), ds)

	request, err := xml.MarshalIndent(v, "", "  ")

	log.Debug("REQUEST: ", string(request))

	if err != nil {
		log.Error("WITH ERROR: ", err)
		return result, err
	}

	response, err := MsoapCall(Runconfig.GetString("AgreedRoutesUrl"), "http://GKOVD/CBD2/GetAgreedSegments", request, Runconfig.GetInt("load_agreed_timeout"))
	if err != nil {
		return result, err
	}

	err = checkXMLResponseForError(response)

	if err != nil {
		log.Error("RESPONSE:\n", string(response), "\nWITH ERROR", err)
		return result, err
	}

	var env GCEnvelop
	err = xml.Unmarshal(response, &env)

	if err != nil {
		log.Error("xml.Unmarshal ERROR", err)
		return result, err
	}

	result = env.Body.Response.Result.Rowset

	return result, nil
}

/*
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cbd2="http://GKOVD/CBD2">
   <soapenv:Header>
      <cbd2:AuthData>
         <!--Optional:-->
         <cbd2:login>MonitorSoft</cbd2:login>
         <!--Optional:-->
         <cbd2:password>123</cbd2:password>
      </cbd2:AuthData>
   </soapenv:Header>
   <soapenv:Body>
      <cbd2:SendXML>
         <!--Optional:-->
         <cbd2:data>
            <trafficReport xmlns="http://GKOVD/CBD2" reportDate="2020-06-16" aftnAddress="УЕЕЕЗДЗЬ">
            ...
*/

type GCSendXML struct {
	XMLName xml.Name `xml:"cbd2:SendXML"`
	Data    []byte   `xml:"cbd2:data"`
}

func GC_SendSvodka(data []byte) ([]byte, error) {

	log := logrus.WithFields(logrus.Fields{
		"file":   "gcprotocol.go",
		"method": "GC_SendSvodka",
	})

	ds := GCSendXML{}
	//ds.Data = data

	v := GC_CreateRequest(Runconfig.GetString("SoapLogin"), Runconfig.GetString("SoapPassword"), ds)

	request, err := xml.MarshalIndent(v, "", "  ")
	// replace  <cbd2:data></cbd2:data> with raw data []byte

	ap_data := []byte("<cbd2:data>")
	ap_data = append(ap_data, data...)
	ap_data = append(ap_data, []byte("</cbd2:data>")...)

	request = bytes.ReplaceAll(request, []byte("<cbd2:data></cbd2:data>"), ap_data)

	log.Debug("REQUEST: ", string(request))

	if err != nil {
		log.Error("WITH ERROR: ", err)
		return []byte{}, err
	}

	response, err := MsoapCall(Runconfig.GetString("svodkaurl"), "http://GKOVD/CBD2/SendXML", request, Runconfig.GetInt("send_xml_timeout"))
	if err != nil {
		return response, err
	}

	err = checkXMLResponseForError(response)

	if err != nil {
		log.Error("RESPONSE:\n", string(response), "\nWITH ERROR", err)
		return response, err
	}
	return response, nil
}
