package odahu.core


test_ds_get_info {
  allow with data.odahu.mapper as {"action": "GET", "resource": "/api/info", "user_roles": ["data_scientist"]}
}
test_ds_post_predict{
  allow with data.odahu.mapper as {"action": "POST", "resource": "/api/predict", "user_roles": ["data_scientist"]}
}
test_admin_get_info {
  allow with data.odahu.mapper as {"action": "GET", "resource": "/api/info", "user_roles": ["admin"]}
}
test_admin_post_predict{
  allow with data.odahu.mapper as {"action": "POST", "resource": "/api/predict", "user_roles": ["admin"]}
}
test_raw_role_get_info{
  allow with data.odahu.mapper as {"action": "GET", "resource": "/api/info", "raw_roles": ["{{.Role}}"]}
}
test_raw_role_post_predict{
  allow with data.odahu.mapper as {"action": "POST", "resource": "/api/predict", "raw_roles": ["{{.Role}}"]}
}
