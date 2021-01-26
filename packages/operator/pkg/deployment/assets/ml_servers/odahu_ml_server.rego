package odahu.core

import data.odahu.mapper
import data.odahu.roles

allow {
  mapper.action == "GET"
  mapper.resource == "/healthz"
}

allow {
  mapper.action == "GET"
  mapper.resource == "/metrics"
}

allow {
  mapper.action == "GET"
  mapper.resource == "/healthcheck"
}

# Raw role
allow {
  mapper.action == ["GET", "POST"][_]
  mapper.resource == ["/api/model/info", "/api/model/invoke"][_]

  mapper.raw_roles[_] == "{{.Role}}"
}

# Fixed roles
allow {
  mapper.action == ["GET", "POST"][_]
  mapper.resource == ["/api/model/info", "/api/model/invoke"][_]

  mapper.user_roles[_] == [roles.data_scientist, roles.admin][_]
}
