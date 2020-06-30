Feature: CORS User

  Scenario: User login CORS
    When I send a "OPTIONS" request on "/users/login"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: Users signup CORS
    When I send a "OPTIONS" request on "/users/login"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

  Scenario: Users me CORS
    When I send a "OPTIONS" request on "/me"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "GET,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"
