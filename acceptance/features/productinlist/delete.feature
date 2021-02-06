@product-in-list
Feature: Add a product to a list

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    When I send a "DELETE" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    Then the status code is 204
    And I send a "GET" request on "/listing/lists/{{ .ListID }}"
    # And the returned products from list are
    #   | ID                                   | Name             | Category name |

  Scenario: KO - List doesn't belong to user
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    And I authenticate with a random user
    When I send a "DELETE" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    Then the status code is 403
    And "$.error" has string value "forbidden-error"

  Scenario: KO - List id is not an UUID
    Given I authenticate with a random user
    When I send a "DELETE" request on "/listing/lists/a/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    # And "$.fields[0].name" has string value "listID"
    # And "$.fields[0].error" has string value "uuid"

  Scenario: KO - Product id is not an UUID
    Given I authenticate with a random user
    When I send a "DELETE" request on "/listing/lists/00000000-0000-0000-0000-000000000000/products/a"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    # And "$.fields[0].name" has string value "productID"
    # And "$.fields[0].error" has string value "uuid"
