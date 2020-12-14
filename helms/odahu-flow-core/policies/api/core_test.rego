package odahu.core

test_swagger {
  allow with data.odahu.mapper.parsed_input as {"action": "GET", "resource": "/swagger/index.html"}
  allow with data.odahu.mapper.parsed_input as {"action": "GET", "resource": "/swagger/"}
  allow with data.odahu.mapper.parsed_input as {"action": "GET", "resource": "/swagger/data.json"}
  allow with data.odahu.mapper.parsed_input as {"action": "GET", "resource": "/swagger"}
  not allow with data.odahu.mapper.parsed_input as {"action": "GET", "resource": "/swagge1r"}
}
