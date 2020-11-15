Feature: CORS List

  Scenario: Lists CORS
    When I send a "OPTIONS" request on "/lists"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,GET,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: List with id CORS
    When I send a "OPTIONS" request on "/lists/hello"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "GET,DELETE,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: Archive a list CORS
    When I send a "OPTIONS" request on "/lists/hello/archive"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"
