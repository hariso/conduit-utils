package main

type CreatePipelineRequest struct {
	Config *Pipeline_Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
}

type Pipeline_Config struct {
	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

type Pipeline struct {
	Id     string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	State  *Pipeline_State  `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	Config *Pipeline_Config `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	// -- children connections ---------------------------------------------------
	ConnectorIds []string `protobuf:"bytes,5,rep,name=connector_ids,json=connectorIds,proto3" json:"connector_ids,omitempty"`
	ProcessorIds []string `protobuf:"bytes,6,rep,name=processor_ids,json=processorIds,proto3" json:"processor_ids,omitempty"`
}

type Pipeline_State struct {
	Status Pipeline_Status `protobuf:"varint,1,opt,name=status,proto3,enum=api.v1.Pipeline_Status" json:"status,omitempty"`
	Error  string          `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

type Pipeline_Status string

const (
	Pipeline_STATUS_RUNNING  Pipeline_Status = "STATUS_RUNNING"
	Pipeline_STATUS_STOPPED  Pipeline_Status = "STATUS_STOPPED"
	Pipeline_STATUS_DEGRADED Pipeline_Status = "STATUS_DEGRADED"
)

type CreateConnectorRequest struct {
	Type       Connector_Type    `protobuf:"varint,1,opt,name=type,proto3,enum=api.v1.Connector_Type" json:"type,omitempty"`
	Plugin     string            `protobuf:"bytes,2,opt,name=plugin,proto3" json:"plugin,omitempty"`
	PipelineId string            `protobuf:"bytes,3,opt,name=pipeline_id,json=pipelineId,proto3" json:"pipeline_id,omitempty"`
	Config     *Connector_Config `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
}

type Connector_Type string

const (
	Connector_TYPE_SOURCE      Connector_Type = "TYPE_SOURCE"
	Connector_TYPE_DESTINATION Connector_Type = "TYPE_DESTINATION"
)

type Connector_Config struct {
	Name     string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Settings map[string]string `protobuf:"bytes,2,rep,name=settings,proto3" json:"settings,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type Connector struct {
	Id         string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Config     *Connector_Config `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
	Type       Connector_Type    `protobuf:"varint,5,opt,name=type,proto3,enum=api.v1.Connector_Type" json:"type,omitempty"`
	Plugin     string            `protobuf:"bytes,6,opt,name=plugin,proto3" json:"plugin,omitempty"`
	PipelineId string            `protobuf:"bytes,7,opt,name=pipeline_id,json=pipelineId,proto3" json:"pipeline_id,omitempty"`
	// -- children connections ---------------------------------------------------
	ProcessorIds []string `protobuf:"bytes,8,rep,name=processor_ids,json=processorIds,proto3" json:"processor_ids,omitempty"`
}

// Type shows the processor type.
type Processor_Type int32

const (
	Processor_TYPE_UNSPECIFIED Processor_Type = 0
	// Processor is a transform.
	Processor_TYPE_TRANSFORM Processor_Type = 1
	// Processor is a filter.
	Processor_TYPE_FILTER Processor_Type = 2
)

// Type shows the processor's parent type.
type Processor_Parent_Type string

const (
	Processor_Parent_TYPE_UNSPECIFIED string = ""
	// Processor parent is a connector.
	Processor_Parent_TYPE_CONNECTOR string = "TYPE_CONNECTOR"
	// Processor parent is a pipeline.
	Processor_Parent_TYPE_PIPELINE string = "TYPE_PIPELINE"
)

type Processor_Parent struct {
	Type string `protobuf:"varint,1,opt,name=type,proto3,enum=api.v1.Processor_Parent_Type" json:"type,omitempty"`
	Id   string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

type Processor_Config struct {
	Settings map[string]string `protobuf:"bytes,1,rep,name=settings,proto3" json:"settings,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type CreateProcessorRequest struct {
	Name   string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type   Processor_Type    `protobuf:"varint,2,opt,name=type,proto3,enum=api.v1.Processor_Type" json:"type,omitempty"`
	Parent *Processor_Parent `protobuf:"bytes,3,opt,name=parent,proto3" json:"parent,omitempty"`
	Config *Processor_Config `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
}
