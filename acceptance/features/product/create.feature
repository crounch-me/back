Feature: Create Product

  Scenario: OK
    Given I create and authenticate with a random user
    And I create these products
      | name             |
      | Creation produit |
    Then the status code is 201
    And "$.name" is a string equal to "Creation produit"
    And "$.id" is a non empty string
