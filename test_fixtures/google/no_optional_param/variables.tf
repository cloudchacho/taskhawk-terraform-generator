{{tfDoNotEditStamp}}

variable "dataflow_tmp_gcs_location" {}

variable "dataflow_template_pubsub_to_pubsub_gcs_path" {}

variable "dataflow_template_pubsub_to_storage_gcs_path" {}

variable "dataflow_zone" {}

variable "enable_firehose_all_messages" {
  default = "false"
}

variable "dataflow_output_directory" {}
