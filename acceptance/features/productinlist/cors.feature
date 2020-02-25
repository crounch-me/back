Feature: CORS Product in list

  Scenario: Product in list CORS
    When I send a "OPTIONS" request on "/lists/1/products/2"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"
