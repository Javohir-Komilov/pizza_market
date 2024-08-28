# Pizzeria Website

This is a simple web application for a pizzeria built as a learning project. The website allows users to browse a menu of pizzas, add pizzas to a shopping cart, register and log in, and place orders. Admin users can manage the menu and order statuses.

## Features

- **Pizza Menu**: Browse a selection of pizzas available for purchase.
- **Shopping Cart**: Add pizzas to your cart and review your selections before placing an order.
- **User Authentication**: Register for an account or log in to place orders.
- **Admin Interface**: Log in with the admin credentials to:
  - Create new pizzas for the menu.
  - Edit existing menu items.
  - Change the status of orders.

## Tech Stack

- **Golang**: The backend is built using Golang, providing robust and efficient server-side logic.
- **templ**: Templating engine used to render HTML pages.
- **Tailwind CSS**: Styling of the website is done with Tailwind CSS.
- **SQLite3**: Database to store user data, menu items, orders, and more.
- **sqlc**: SQL code generation tool used to interact with the SQLite3 database.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Javohir-Komilov/pizzeria-website.git
   cd pizzeria-website
