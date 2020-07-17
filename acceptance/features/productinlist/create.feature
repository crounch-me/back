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
          "productId": "{{ .ProductID }}",
          "listId": "{{ .ListID }}"
        }
      """
    When I send a "POST" request on "/lists/{{ .ListID }}/products/{{ .ProductID}}"
    Then the status code is 201
    And "$.productId" is a string equal to "{{ .ProductID }}"
    And "$.listId" is a string equal to "{{ .ListID }}"

  Scenario: KO - Product already in list
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
          "productId": "{{ .ProductID }}",
          "listId": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/lists/{{ .ListID }}/products/{{ .ProductID}}"
    When I send a "POST" request on "/lists/{{ .ListID }}/products/{{ .ProductID}}"
    Then the status code is 409
    And "$.error" has string value "duplicate-product-in-list-error"

  Scenario: KO - List id is not an UUID
    Given I authenticate with a random user
    When I send a "POST" request on "/lists/a/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "listID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - Product id is not an UUID
    Given I authenticate with a random user
    When I send a "POST" request on "/lists/00000000-0000-0000-0000-000000000000/products/a"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "productID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - List not found
    Given I authenticate with a random user
    When I send a "POST" request on "/lists/00000000-0000-0000-0000-000000000000/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" has string value "list-not-found-error"

  Scenario: KO - Product not found
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "POST" request on "/lists/{{ .ListID }}/products/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" has string value "product-not-found-error"
