package tv_return_service

import (
	"TVTestApp/dbconn"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type TvXml struct {
	XMLName xml.Name `xml:"tvs"`
	Tvs     []TV     `xml:"tv"`
}

type TV struct {
	XMLName xml.Name `xml:"tv"`
	ID      int64    `xml:"id"`
	Count   int      `xml:"count"`
	Readed  bool     `xml:"readed,attr"`
}

const fileName = "tv_returns.xml"

func ReadXML() (interface{}, error) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var TvXmlInfo TvXml
	err = xml.Unmarshal(byteValue, &TvXmlInfo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return TvXmlInfo, err
}

func WriteData(TvXmlInfo TvXml) error {
	for index, TvXml := range TvXmlInfo.Tvs {
		if TvXml.Readed != true {
			TV, err := dbconn.GetTv(TvXml.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}
			count := TV.Count + TvXml.Count
			if count < 0 {
				fmt.Println("return data is incorrect")
				return err
			}
			fmt.Printf("Executing return tvs. ID:%v, Count:%v, OldCount:%v\n", TV.ID, count, TV.Count)
			err = dbconn.UpdateTvsCount(TV.ID, count)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Printf("Successful execute tvs. ID:%v, Count:%v, OldCount:%v\n", TV.ID, count, TV.Count)
			TvXmlInfo.Tvs[index].Readed = true
		}
		modifiedXml, err := xml.Marshal(TvXmlInfo)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = ioutil.WriteFile(fileName, modifiedXml, 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
