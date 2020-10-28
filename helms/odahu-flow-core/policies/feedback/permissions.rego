package odahu.permissions

import data.odahu.roles

permissions := {
	roles.data_scientist: [
    	[".*", "^/api/v1/feedback.*"],
    ],
  roles.admin : [
      [".*", ".*"]
  ]
}