package odahu.mapper

import data.odahu.roles

roles_map = {
	"odahu_admin": roles.admin,
  "odahu_data_scientist": roles.data_scientist,
  "odahu_viewer": roles.viewer
}

jwt = input.attributes.metadata_context.filter_metadata["envoy.filters.http.jwt_authn"].fields.jwt_payload

raw_roles[role]{
	role = jwt.Kind.StructValue.fields.realm_access.Kind.StructValue.fields.roles.Kind.ListValue.values[_].Kind.StringValue
}

user_roles[role]{
	role = roles_map[raw_roles[_]]
}

parsed_input = {
  "action": input.attributes.request.http.method,
  "resource": input.attributes.request.http.path,
  "user": {
    "roles": user_roles
  }
}