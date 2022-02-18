const dataset = {
  "namespace": {
    "name": "default",
    "description": "Default namespace",
    "object_store" :{
      "name":"default",
      "description": "Default object store",
      "endpoint": "http://localhost",
      "type":"minio"
    },
    "bucket_name": "data"
  },
  "name":"my",
  "description":"testing this out",
  "created":"2021-06-09T18:21:41.179278Z",
  "root_directory":"my/",
  "owner": {
    "username":"admin",
    "memberships":null
  },
  "permissions":[
    {
      "group": {
        "group_name":"admin-hoss-default-group"
      },
      "permission": "rw"
    },
    {
      "group": {
        "group_name": "group"
      },
      "permission": "rw"
    },
    {
      "group": {
        "group_name":"sdds"
      },
      "permission": "r"
    },
    {
      "group": {
        "group_name": "my-group"
      },
      "permission": "r"
    },
    {
      "group": {
        "group_name":"your-group"
      },
      "permission":"rw"
    }
  ]
};


export default dataset;
