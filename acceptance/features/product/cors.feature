Feature: CORS Product

  Scenario: Products CORS
    When I send a "OPTIONS" request on "/products"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"
