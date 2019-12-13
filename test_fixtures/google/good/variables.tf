{{tfDoNotEditStamp}}

variable "dataflow_tmp_gcs_location" {
  default = "gs://myBucket/tmp"
}

variable "dataflow_template_pubsub_to_pubsub_gcs_path" {
  default = "gs://dataflow-templates/2019-04-03-00/Cloud_PubSub_to_Cloud_PubSub"
}

variable "dataflow_template_pubsub_to_storage_gcs_path" {
  default = "gs://dataflow-templates/2019-04-03-00/Cloud_PubSub_to_GCS_Text"
}

variable "dataflow_zone" {
  default = "us-west2-a"
}

variable "dataflow_region" {
  default = "us-west1"
}

variable "dataflow_output_directory" {
  default = "gs://myBucket/taskhawkBackup/"
}

variable "enable_firehose_all_messages" {
  default = "true"
}

variable "alerting_project" {
  default = "alerting-project"
}
