package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

type Field struct {
	Type     string `json:"type,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Field    string `json:"field,omitempty"`
}

type Schema struct {
	Type     string  `json:"type,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
	Optional bool    `json:"optional,omitempty"`
	Name     string  `json:"name,omitempty"`
}

func toString(schema Schema) string {
	bytes, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

var (
	serverCerPem = "Bag Attributes\n    friendlyName: caroot\n    2.16.840.1.113894.746875.1.1: <Unsupported tag 6>\nsubject=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\nissuer=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\n-----BEGIN CERTIFICATE-----\nMIIDHjCCAgYCCQC9iilqJUAoxzANBgkqhkiG9w0BAQsFADBRMRIwEAYDVQQDDAls\nb2NhbGhvc3QxDDAKBgNVBAsMA0NJQTEMMAoGA1UECgwDUkVBMRIwEAYDVQQHDAlN\nZWxib3VybmUxCzAJBgNVBAYTAkFVMB4XDTIyMDMyMTEyNDcyNloXDTQ5MDgwNTEy\nNDcyNlowUTESMBAGA1UEAwwJbG9jYWxob3N0MQwwCgYDVQQLDANDSUExDDAKBgNV\nBAoMA1JFQTESMBAGA1UEBwwJTWVsYm91cm5lMQswCQYDVQQGEwJBVTCCASIwDQYJ\nKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK8eUXQWIEvqH6gWg9CyU1FZ9I5ZfeSz\n5BgEwQelG3YRiBmN4MXQmVErvy8JEdC9AbDdNvwsWlBD1xoWC0S2Q2qMhF6M03ny\nrrx0OwKNxdNwZvrMCin6adVS66x4R83X/YprZiS0fMtZHnrPsEVZxw7QSObGPnUV\nqinVqZh4Mo2N7tbxYa6ZALXgDf0yXbzGOGENuEaw9+5H01+6wDAwoxmm3pgQ0bF1\nyrqGh6P5ePtbaI6C+WBW/u0HgXUgyJaQA3vZIS6cOnwf76osPkFiCp5LtjTWblBd\nBwEmjtMu4n6/QEsUNM93lP/iJ7y8LWGrMFxL1700KFWlkJ8CIbpdjH0CAwEAATAN\nBgkqhkiG9w0BAQsFAAOCAQEAAJXIy4onOdKSceG3gjJHjCIZf74/Ka7rIhCnVyn1\n/EiTub70Df1b2vsPl25axv3ujMETubUzSzyyRSQ/5o61T5nZXPmmbTtU/7f1lLWk\n6tiHyKRCPd6InLnGFqhjmKp052LCwX0ZqUy2x+/uXxrYI7+wcB9QpDQsw4nnfDhb\nkym8hzUIu6IST1TbHFD7NoTA1L/qVlVQ8Sj6XVGQKmwBPxX/i1wHFTDrnS4uskvK\nFkPUsd6OmYbguwS3Ktj40C/pV3Z5OS/kR4+pO349I+b42ImzWxpgMRSVuI4y0Lvk\nm6GdbnJPvzqZT5yZmv05j6LQXkm5ugqPmOARMKrSrWACUA==\n-----END CERTIFICATE-----\n"
	clientKeyPem = "Bag Attributes\n    friendlyName: broker\n    localKeyID: 54 69 6D 65 20 31 36 34 37 38 36 37 32 30 38 36 35 35 \nKey Attributes: <No Attributes>\n-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCQTuVvpsn8Eu3b\nerPVw576+bifl9QJTtDdaYHeywDRmrp3XO7VUOsRf6Wacm+4+2Uvbi0NVrhWbLwe\n7/jNoOhgDXCS85A6FnZU+LBfwJBLBJjPFZ654rFMjz+kmuHf2b5J7LkqXAuSTAFe\nW11EiMchiREg9inhoxwGt2qviEUYabZLJR31pr2qd7jGiLK9EPY3y8UaW/15HUXq\nQhWVntAQPFTtitnue7fet4OwcY8nCkJ5yOnxRznX78k877ycfgKLJwHkETwHWp2x\nXwJyTz0mJ0nhbYadacGHieKgqu1zUC9f7Tz4AUOB+coEALqKnRPz9nKmT4lUQwD+\nzOr7tmMxAgMBAAECggEAabZyGuV6577SIcr0PG8OYlpXJgoqGRt0pA3rRlM96U5I\ntLIOf5PEb9Ard0XHlCINUL6MIE5bwWvsL1mp0LDEKcEOq4fjKrpTuxFm2u4Mhff7\nHRCAczmemjAB9kpDlyFCZZMVXfOJwoUNJ5sUauUrwuRO+O97ZMCBAmaQr7/KpgOI\n6OxPYv0WE6uqaijFWp7HnXc5Xcw0zo5riEh+5ZxjaehJFWoFpOKpOgL+pkMqfAZb\nVGTY8BxRJA11nR/Do34kz8WahX6Us3DCMpMMWJcg8ss80JRFyj30o683tOfBeTdP\nHVJrRz1GdW7wHnGR94D/a3TuUISJ1U++GUMdXmQNFQKBgQDHrgG2wwzGQ8Uv+cWD\nLj02flAA6Xrln4JQpelrQZrrYcYbPevTordeeXnY+COtGF04ix5dQWsWAgYlRSnr\n1s56en+RsnhvNggoHz48nPXNwbp/XZx1or8fJUMwLicwosFtLF7DxLMF4tq1pdkX\nZrr9UJO7LtirgAOqiZzNsd9xOwKBgQC5AsFh7CHSwNo/x6G6XJ9c4hL0dYeuxbbo\n4a08NageTRWXxBdpkYn4YBrLk49jA3JnDrfJuJfJ7jQdGBLwNdTFWZp+wB/QT2su\nkE2Ur+ysid95xLM1aCw17RfzpSXhkhMUJTRSYgIyqkkClRnRqhUKDowyt0AfdzCW\nGAPTdac2gwKBgCquDr+5wSk/ow42HPmFEKBtLzyCqzoZdgk27UV3qF1XcLix6444\n4WjYHis6HqYI5yQG2F6mdPUnSZj9x5AZQdj8BfhmZUegDO5Gf08FXaS1G9/Nanva\nZW+Kz2mk88t5fk6PhVHi4UEI1CavZE+ULbOnXWxM/xLpMd9pupJcyp2xAoGALs+Z\nqnMao76T+itCqmqhD9lLvnq2V+xCuW3QbTmOTgxm+D1vRxDB/gwi+3tcfkry+Uxq\nCCoijb8thGcA87JLIZvoUUW/Ru+xSNjOKF7S3V0NJDw2s76l4QcaVlVk3kwdc61u\nLaIKuFMJohOjsr78D81af8KKAOwhaPiujyRnqI0CgYAmyjg1o85Od/Q0MwM2zQ9w\ng14mYRjaLCvvYdCe62tw9B8JTfO5n6kIFokCUJY6lsZbjE4su2sES5VftoizUk2p\njGoznm1taoeEby59583cf79ijBWj+gt8KZl3KvM1vjN55lZ4mUS8yIDzzlShoXOf\nexA8AP+RDi6/PwW4SYQpmg==\n-----END PRIVATE KEY-----\n"
	clientCerPem = "Bag Attributes\n    friendlyName: broker\n    localKeyID: 54 69 6D 65 20 31 36 34 37 38 36 37 32 30 38 36 35 35 \nsubject=/C=AU/ST=VIC/L=Melbourne/O=REA/OU=CIA/CN=localhost\nissuer=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\n-----BEGIN CERTIFICATE-----\nMIIDLDCCAhQCCQCGHfkQtzwLDTANBgkqhkiG9w0BAQUFADBRMRIwEAYDVQQDDAls\nb2NhbGhvc3QxDDAKBgNVBAsMA0NJQTEMMAoGA1UECgwDUkVBMRIwEAYDVQQHDAlN\nZWxib3VybmUxCzAJBgNVBAYTAkFVMB4XDTIyMDMyMTEyNDcyN1oXDTQ5MDgwNTEy\nNDcyN1owXzELMAkGA1UEBhMCQVUxDDAKBgNVBAgTA1ZJQzESMBAGA1UEBxMJTWVs\nYm91cm5lMQwwCgYDVQQKEwNSRUExDDAKBgNVBAsTA0NJQTESMBAGA1UEAxMJbG9j\nYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkE7lb6bJ/BLt\n23qz1cOe+vm4n5fUCU7Q3WmB3ssA0Zq6d1zu1VDrEX+lmnJvuPtlL24tDVa4Vmy8\nHu/4zaDoYA1wkvOQOhZ2VPiwX8CQSwSYzxWeueKxTI8/pJrh39m+Sey5KlwLkkwB\nXltdRIjHIYkRIPYp4aMcBrdqr4hFGGm2SyUd9aa9qne4xoiyvRD2N8vFGlv9eR1F\n6kIVlZ7QEDxU7YrZ7nu33reDsHGPJwpCecjp8Uc51+/JPO+8nH4CiycB5BE8B1qd\nsV8Cck89JidJ4W2GnWnBh4nioKrtc1AvX+08+AFDgfnKBAC6ip0T8/Zypk+JVEMA\n/szq+7ZjMQIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQAh1sobuLh1uN6qZJOvV6vS\nG2a182VhX0ktBxXTZyXshSKa0lT93vtPgMEz/xRQ3H6ZVdEh6+GbY7jLIYjiqFhm\neuDnb6A+SP1VdosSPY6pg9tNWVIcrVTeUltrbJGYp7HTyaAvgqc5fzinhEmvbwxr\n/A3LBUjr/WrwzTCq/lwhQwjE61EET9SV/fzcF+I8I8SF5uVwHEnlyV9FaFYg36Ba\nZBt3Mfvf1Ai477N3npv7j9OoenCTxuIr0jrQJ7QT1pq1wEjckSAbuzpqqg0TXy1l\nXBkCk5xXHRriCnPBtL6VMlnLwpxap1nOTHWkp70z8HerW44+GfnoMI67G5Q16GBX\n-----END CERTIFICATE-----\nBag Attributes\n    friendlyName: C=AU,L=Melbourne,O=REA,OU=CIA,CN=localhost\nsubject=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\nissuer=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\n-----BEGIN CERTIFICATE-----\nMIIDHjCCAgYCCQC9iilqJUAoxzANBgkqhkiG9w0BAQsFADBRMRIwEAYDVQQDDAls\nb2NhbGhvc3QxDDAKBgNVBAsMA0NJQTEMMAoGA1UECgwDUkVBMRIwEAYDVQQHDAlN\nZWxib3VybmUxCzAJBgNVBAYTAkFVMB4XDTIyMDMyMTEyNDcyNloXDTQ5MDgwNTEy\nNDcyNlowUTESMBAGA1UEAwwJbG9jYWxob3N0MQwwCgYDVQQLDANDSUExDDAKBgNV\nBAoMA1JFQTESMBAGA1UEBwwJTWVsYm91cm5lMQswCQYDVQQGEwJBVTCCASIwDQYJ\nKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK8eUXQWIEvqH6gWg9CyU1FZ9I5ZfeSz\n5BgEwQelG3YRiBmN4MXQmVErvy8JEdC9AbDdNvwsWlBD1xoWC0S2Q2qMhF6M03ny\nrrx0OwKNxdNwZvrMCin6adVS66x4R83X/YprZiS0fMtZHnrPsEVZxw7QSObGPnUV\nqinVqZh4Mo2N7tbxYa6ZALXgDf0yXbzGOGENuEaw9+5H01+6wDAwoxmm3pgQ0bF1\nyrqGh6P5ePtbaI6C+WBW/u0HgXUgyJaQA3vZIS6cOnwf76osPkFiCp5LtjTWblBd\nBwEmjtMu4n6/QEsUNM93lP/iJ7y8LWGrMFxL1700KFWlkJ8CIbpdjH0CAwEAATAN\nBgkqhkiG9w0BAQsFAAOCAQEAAJXIy4onOdKSceG3gjJHjCIZf74/Ka7rIhCnVyn1\n/EiTub70Df1b2vsPl25axv3ujMETubUzSzyyRSQ/5o61T5nZXPmmbTtU/7f1lLWk\n6tiHyKRCPd6InLnGFqhjmKp052LCwX0ZqUy2x+/uXxrYI7+wcB9QpDQsw4nnfDhb\nkym8hzUIu6IST1TbHFD7NoTA1L/qVlVQ8Sj6XVGQKmwBPxX/i1wHFTDrnS4uskvK\nFkPUsd6OmYbguwS3Ktj40C/pV3Z5OS/kR4+pO349I+b42ImzWxpgMRSVuI4y0Lvk\nm6GdbnJPvzqZT5yZmv05j6LQXkm5ugqPmOARMKrSrWACUA==\n-----END CERTIFICATE-----\nBag Attributes\n    friendlyName: caroot\n    2.16.840.1.113894.746875.1.1: <Unsupported tag 6>\nsubject=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\nissuer=/CN=localhost/OU=CIA/O=REA/L=Melbourne/C=AU\n-----BEGIN CERTIFICATE-----\nMIIDHjCCAgYCCQC9iilqJUAoxzANBgkqhkiG9w0BAQsFADBRMRIwEAYDVQQDDAls\nb2NhbGhvc3QxDDAKBgNVBAsMA0NJQTEMMAoGA1UECgwDUkVBMRIwEAYDVQQHDAlN\nZWxib3VybmUxCzAJBgNVBAYTAkFVMB4XDTIyMDMyMTEyNDcyNloXDTQ5MDgwNTEy\nNDcyNlowUTESMBAGA1UEAwwJbG9jYWxob3N0MQwwCgYDVQQLDANDSUExDDAKBgNV\nBAoMA1JFQTESMBAGA1UEBwwJTWVsYm91cm5lMQswCQYDVQQGEwJBVTCCASIwDQYJ\nKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK8eUXQWIEvqH6gWg9CyU1FZ9I5ZfeSz\n5BgEwQelG3YRiBmN4MXQmVErvy8JEdC9AbDdNvwsWlBD1xoWC0S2Q2qMhF6M03ny\nrrx0OwKNxdNwZvrMCin6adVS66x4R83X/YprZiS0fMtZHnrPsEVZxw7QSObGPnUV\nqinVqZh4Mo2N7tbxYa6ZALXgDf0yXbzGOGENuEaw9+5H01+6wDAwoxmm3pgQ0bF1\nyrqGh6P5ePtbaI6C+WBW/u0HgXUgyJaQA3vZIS6cOnwf76osPkFiCp5LtjTWblBd\nBwEmjtMu4n6/QEsUNM93lP/iJ7y8LWGrMFxL1700KFWlkJ8CIbpdjH0CAwEAATAN\nBgkqhkiG9w0BAQsFAAOCAQEAAJXIy4onOdKSceG3gjJHjCIZf74/Ka7rIhCnVyn1\n/EiTub70Df1b2vsPl25axv3ujMETubUzSzyyRSQ/5o61T5nZXPmmbTtU/7f1lLWk\n6tiHyKRCPd6InLnGFqhjmKp052LCwX0ZqUy2x+/uXxrYI7+wcB9QpDQsw4nnfDhb\nkym8hzUIu6IST1TbHFD7NoTA1L/qVlVQ8Sj6XVGQKmwBPxX/i1wHFTDrnS4uskvK\nFkPUsd6OmYbguwS3Ktj40C/pV3Z5OS/kR4+pO349I+b42ImzWxpgMRSVuI4y0Lvk\nm6GdbnJPvzqZT5yZmv05j6LQXkm5ugqPmOARMKrSrWACUA==\n-----END CERTIFICATE-----\n"

	generatorPlugin      = "pkg/plugins/generator/generator"
	kafkaConnectorPlugin = "/home/haris/projects/conduitio/conduit-kafka-connect-wrapper/dist/conduit-kafka-connect-wrapper"
	kcSourceDebug        = "/home/haris/projects/personal/conduit-utils/kafka-connect-grpc-source.sh"
	kcDestinationDebug   = "/home/haris/projects/personal/conduit-utils/kafka-connect-grpc-destination.sh"
	filePlugin           = "builtin:file"
	kafkaPluginOld       = "pkg/plugins/kafka/kafka"
	kafkaPlugin          = "/home/haris/projects/conduitio/conduit-connector-kafka/conduit-connector-kafka"
	pgPlugin             = "builtin:postgres"

	customersSchema = toString(Schema{
		Type: "struct",
		Fields: []Field{
			{
				Type:     "int32",
				Optional: true,
				Field:    "id",
			},
			{
				Type:     "string",
				Optional: true,
				Field:    "name",
			},
			{
				Type:     "boolean",
				Optional: true,
				Field:    "trial",
			},
		},
		Optional: false,
		Name:     "customers",
	})

	nameOnlySchema = toString(Schema{
		Type:     "string",
		Optional: false,
		Name:     "foobar",
	})

	customersPartialSchema = toString(Schema{
		Type: "struct",
		Fields: []Field{
			{
				Type:     "boolean",
				Optional: true,
				Field:    "joined",
			},
		},
		Optional: false,
		Name:     "customers",
	})

	countriesSchema = toString(Schema{
		Type: "struct",
		Fields: []Field{
			{
				Type:     "int32",
				Optional: true,
				Field:    "id",
			},
			{
				Type:     "string",
				Optional: true,
				Field:    "name",
			},
			{
				Type:     "boolean",
				Optional: true,
				Field:    "large",
			},
		},
		Optional: false,
		Name:     "countries",
	})

	mySQLConnUrl     = "jdbc:mysql://haris:password@localhost:3306/haris-test-db"
	pgUrl            = "postgresql://repmgr:repmgrmeroxa@localhost/conduit-test-db?sslmode=disable"
	s3AccessKey      = "AKIA2AI5UVYKM4GWKKIL"
	s3SecretKey      = "*****"
	s3Bucket         = "hariso"
	openWeatherAppId = "*****"
)

func TestFoo(t *testing.T) {
	fmt.Println(customersPartialSchema)
}

func TestStartPipelineOnly(t *testing.T) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	err := startPipeline(client, "080e2618-87b6-448b-a5ef-9d567842308f")
	if err != nil {
		panic(err)
	}
}

func TestStopPipeline(t *testing.T) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	err := stopPipeline(client, "c76ffb45-426e-4d53-ac4e-ef0890a08467")
	if err != nil {
		panic(err)
	}
}

func TestRecreatePipeline(t *testing.T) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	pipelineId, err := createPipeline(client)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createConnectPgSource(client, pipelineId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = createConnectSnowflakeDestination(client, pipelineId)
	if err != nil {
		t.Fatal(err)
	}

	err = startPipeline(client, pipelineId)
	if err != nil {
		t.Fatal(err)
	}
}

func createConnectSnowflakeDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectSnowflakeDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kcDestinationDebug,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-snowflake-sink",
			Settings: map[string]string{
				"wrapper.connector.class":             "com.snowflake.kafka.connector.SnowflakeSinkConnector",
				"wrapper.schema.autogenerate.enabled": "true",
				"wrapper.schema.autogenerate.name":    "CUSTOMERS_TEST",
				"tasks.max":                           "8",
				"topics":                              "customers",
				"name":                                "mysnowflakesink",
				"snowflake.topic2table.map":           "customers:CUSTOMERS_TEST",
				"buffer.count.records":                "1",
				"buffer.flush.time":                   "0",
				"buffer.size.bytes":                   "1",
				"snowflake.url.name":                  "abc.us-east-2.aws.snowflakecomputing.com:443",
				"snowflake.user.name":                 "meroxa_user",
				"snowflake.private.key":               "********",
				"snowflake.database.name":             "MEROXA_DB",
				"snowflake.schema.name":               "STREAM_DATA",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectMongoDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectMongoDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-mongo-destination",
			Settings: map[string]string{
				"task.class":    "com.mongodb.kafka.connect.sink.MongoSinkTask",
				"connectorName": "my-mongo-destination",
				"pipelineId":    pipelineId,
				"database":      "haris-test-db",
				"collection":    "my-test-collection",
				"topics":        "countries",
				"schema":        countriesSchema,
			},
		},
	}

	return createConnector(client, req)
}

func createConnector(client *http.Client, req CreateConnectorRequest) (string, error) {
	pipelineCreateUrl := "http://localhost:8080/v1/connectors"

	resp := &Connector{}
	err := post(client, pipelineCreateUrl, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func createGeneratorSource(client *http.Client, pipelineId string) (interface{}, error) {
	fmt.Println("createGeneratorSource")
	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     generatorPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-generator-source",
			Settings: map[string]string{
				"recordCount": "5",
				"readTime":    "500ms",
				"fields":      "id:int,name:string,trial:bool",
			},
		},
	}

	return createConnector(client, req)
}

func stopPipeline(client *http.Client, pipelineId string) error {
	fmt.Println("starting pipeline")
	pipelineCreateUrl := fmt.Sprintf("http://localhost:8080/v1/pipelines/%s/stop", pipelineId)

	req, err := http.NewRequest("POST", pipelineCreateUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("got status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}

func startPipeline(client *http.Client, pipelineId string) error {
	fmt.Println("starting pipeline")

	pipelineCreateUrl := fmt.Sprintf("http://localhost:8080/v1/pipelines/%s/start", pipelineId)

	req, err := http.NewRequest("POST", pipelineCreateUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("got status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}

func createConnectPgSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectPgSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-pg-source",
			Settings: map[string]string{
				"wrapper.connector.class":  "io.aiven.connect.jdbc.JdbcSourceConnector",
				"connection.url":           "jdbc:postgresql://localhost/meroxadb",
				"connection.user":          "meroxauser",
				"connection.password":      "meroxapass",
				"tables":                   "customers",
				"mode":                     "incrementing",
				"incrementing.column.name": "id",
				"topic.prefix":             "customers",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectHttpSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectHttpSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-http-source",
			Settings: map[string]string{
				"wrapper.task.class":     "com.github.castorm.kafka.connect.http.HttpSourceTask",
				"wrapper.connector.name": "my-http-source",
				"wrapper.pipeline.id":    pipelineId,

				"http.request.method": "GET",
				"http.request.url":    "https://api.openweathermap.org/data/2.5/weather?lat=43.87&lon=18.42&appid=" + openWeatherAppId,
				"kafka.topic":         "openweather-data-topic",
			},
		},
	}

	return createConnector(client, req)
}

func createDebeziumMySQLSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createDebeziumMySQLSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "debezium-mysql-source",
			Settings: map[string]string{
				"wrapper.task.class":     "io.debezium.connector.mysql.MySqlConnectorTask",
				"wrapper.connector.name": "debezium-mysql-source",
				"wrapper.pipeline.id":    pipelineId,

				"database.hostname":                        "localhost",
				"database.port":                            "3306",
				"database.user":                            "debezium",
				"database.password":                        "dbz",
				"database.server.id":                       "184054",
				"database.server.name":                     "dbserver1",
				"database.include.list":                    "inventory",
				"database.history.kafka.bootstrap.servers": "localhost:9092",
				"database.history.kafka.topic":             "schema-changes.inventory",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectHttpDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectHttpDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-http-sink",
			Settings: map[string]string{
				"wrapper.connector.class": "uk.co.threefi.connect.http.HttpSinkConnector",
				"wrapper.schema":          nameOnlySchema,

				"http.api.url":   "http://localhost:1080/mock-destination",
				"request.method": "POST",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectPgDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectPgDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-pg-sink",
			Settings: map[string]string{
				"wrapper.task.class": "io.aiven.connect.jdbc.sink.JdbcSinkTask",
				"wrapper.schema":     customersSchema,
				// "schema.autogenerate.overrides": customersPartialSchema,
				"wrapper.schema.autogenerate.enabled": "true",
				"wrapper.schema.autogenerate.name":    "customers_copy",
				"wrapper.connector.name":              "my-pg-sink",
				"wrapper.pipeline.id":                 pipelineId,

				"connection.url":      "jdbc:postgresql://localhost/conduit-test-db",
				"connection.user":     "meroxauser",
				"connection.password": "meroxapass",
				"auto.create":         "true",
				"auto.evolve":         "true",
				"batch.size":          "10",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectMySQLSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectMySQLSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-pg-sink",
			Settings: map[string]string{
				"task.class":               "io.aiven.connect.jdbc.source.JdbcSourceTask",
				"connection.url":           mySQLConnUrl,
				"tables":                   "customers",
				"topic.prefix":             "foobar_",
				"mode":                     "incrementing",
				"incrementing.column.name": "id",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectMySQLDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectMySQLDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-mysql-sink",
			Settings: map[string]string{
				"task.class": "io.aiven.connect.jdbc.sink.JdbcSinkTask",
				// "schema":         countriesSchema,
				"schema.autogenerate": "true",
				"connectorName":       "my-mysql-sink",
				"pipelineId":          pipelineId,
				"connection.url":      mySQLConnUrl,
				"auto.create":         "true",
				"auto.evolve":         "true",
				"batch.size":          "10",
			},
		},
	}

	return createConnector(client, req)
}

func createConnectS3Sink(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createConnectS3Sink")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaConnectorPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-s3-sink",
			Settings: map[string]string{
				"task.class":             "io.aiven.kafka.connect.s3.S3SinkTask",
				"aws.access.key.id":      s3AccessKey,
				"aws.secret.access.key":  s3SecretKey,
				"aws.s3.bucket.name":     s3Bucket,
				"format.output.type":     "json",
				"file.compression.type":  "none",
				"format.output.envelope": "false",
				"file.name.template":     "{{key}}",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-kafka-destination",
			Settings: map[string]string{
				"servers": "localhost:9092",
				"topic":   "my-kafka-dest-" + uuid.NewString(),
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaDestinationTLS(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaDestinationTLS")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-kafka-destination",
			Settings: map[string]string{
				"servers":            "localhost:19092",
				"topic":              "my-kafka-sink",
				"clientCert":         clientCerPem,
				"clientKey":          clientKeyPem,
				"caCert":             serverCerPem,
				"insecureSkipVerify": "true",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaDestinationTLS_SASL(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaDestinationTLS_SASL")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-kafka-destination",
			Settings: map[string]string{
				"servers":            "localhost:39092",
				"topic":              "my-kafka-dest",
				"clientCert":         clientCerPem,
				"clientKey":          clientKeyPem,
				"caCert":             serverCerPem,
				"insecureSkipVerify": "true",
				"saslUsername":       "admin",
				"saslPassword":       "admin-secret",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaDestinationTLS_SASL_SCRAM(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaDestinationTLS_SASL_SCRAM")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-kafka-dest",
			Settings: map[string]string{
				"servers":       "localhost:19094",
				"topic":         "my-kafka-dest",
				"saslMechanism": "SCRAM-SHA-256",
				"saslUsername":  "metricsreporter",
				"saslPassword":  "password",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaDestinationSASL(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaDestinationSASL")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-kafka-destination",
			Settings: map[string]string{
				"servers":      "localhost:29092",
				"topic":        "my-kafka-dest",
				"saslUsername": "admin",
				"saslPassword": "admin-secret",
			},
		},
	}

	return createConnector(client, req)
}

func createFileSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createFileSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     filePlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-file-source",
			Settings: map[string]string{
				"path": "/Users/haris/Desktop/file-source.txt",
			},
		},
	}

	return createConnector(client, req)
}

func createFileDestination(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createFileDestination")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     filePlugin,
		Type:       Connector_TYPE_DESTINATION,
		Config: &Connector_Config{
			Name: "my-file-destination",
			Settings: map[string]string{
				"path": "/Users/haris/Desktop/file-destination.txt",
			},
		},
	}

	return createConnector(client, req)
}

func createPgSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createPgSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     pgPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-pg-source",
			Settings: map[string]string{
				"table": "customers",
				"url":   pgUrl,
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaSource(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaSource")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-pg-source",
			Settings: map[string]string{
				"servers":           "localhost:9092",
				"topic":             "my-kafka-source",
				"readFromBeginning": "true",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaSourceTLS(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaSourceTLS")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-kafka-source",
			Settings: map[string]string{
				"servers":            "localhost:19092",
				"topic":              "my-kafka-source",
				"readFromBeginning":  "true",
				"clientCert":         clientCerPem,
				"clientKey":          clientKeyPem,
				"caCert":             serverCerPem,
				"insecureSkipVerify": "true",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaSourceTLS_SASL(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaSourceTLS_SASL")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-kafka-source",
			Settings: map[string]string{
				"servers":            "localhost:39092",
				"topic":              "my-kafka-source",
				"readFromBeginning":  "true",
				"clientCert":         clientCerPem,
				"clientKey":          clientKeyPem,
				"caCert":             serverCerPem,
				"insecureSkipVerify": "true",
				"saslUsername":       "admin",
				"saslPassword":       "admin-secret",
			},
		},
	}

	return createConnector(client, req)
}
func createKafkaSourceSASL(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaSourceSASL")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-kafka-source",
			Settings: map[string]string{
				"servers":           "localhost:29092",
				"topic":             "my-kafka-source",
				"readFromBeginning": "true",
				"saslUsername":      "admin",
				"saslPassword":      "admin-secret",
			},
		},
	}

	return createConnector(client, req)
}

func createKafkaSource_SASL_SCRAM(client *http.Client, pipelineId string) (string, error) {
	fmt.Println("createKafkaSource_SASL_SCRAM")

	req := CreateConnectorRequest{
		PipelineId: pipelineId,
		Plugin:     kafkaPlugin,
		Type:       Connector_TYPE_SOURCE,
		Config: &Connector_Config{
			Name: "my-kafka-source",
			Settings: map[string]string{
				"servers":           "localhost:19094",
				"topic":             "my-kafka-source",
				"readFromBeginning": "true",
				"saslMechanism":     "SCRAM-SHA-256",
				"saslUsername":      "metricsreporter",
				"saslPassword":      "password",
			},
		},
	}

	return createConnector(client, req)
}

func createPipeline(client *http.Client) (string, error) {
	fmt.Println("creating pipeline")
	pipelineCreateUrl := "http://localhost:8080/v1/pipelines"

	reqBody := CreatePipelineRequest{
		Config: &Pipeline_Config{
			Description: "My new pipeline",
			Name:        fmt.Sprintf("pipeline-name-%v", uuid.New().String()),
		},
	}

	resp := &Pipeline{}
	err := post(client, pipelineCreateUrl, reqBody, resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func post(client *http.Client, url string, reqBody interface{}, parseInto interface{}) error {
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("got status code %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, &parseInto)
	if err != nil {
		return err
	}
	return nil
}
