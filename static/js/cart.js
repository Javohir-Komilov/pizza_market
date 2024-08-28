function updateQuantity(itemId, delta) {
    const quantityElement = document.querySelector(`#quantity-${itemId}`);
    const subtotalElement = document.querySelector(`#subtotal-${itemId}`);
    const priceElement = document.querySelector(`#price-${itemId}`);
    let quantity = parseInt(quantityElement.textContent) + delta;
    
    if (quantity < 1) {
        quantity = 1;
    }

    fetch('/user/cart/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ item_id: itemId, quantity: quantity }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            quantityElement.textContent = quantity;
            const price = parseFloat(priceElement.textContent);
            const newSubtotal = (price * quantity).toFixed(2);
            subtotalElement.textContent = newSubtotal;
            updateCartTotal();
        }
    });
}

function removeItem(itemId) {
    fetch('/user/cart/remove', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ item_id: itemId }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            const itemElement = document.querySelector(`#cart-item-${itemId}`);
            itemElement.remove();
            updateCartTotal();
        }
    });
}

function updateCartTotal() {
    const subtotals = document.querySelectorAll('[id^="subtotal-"]');
    let total = 0;
    subtotals.forEach(subtotal => {
        total += parseFloat(subtotal.textContent);
    });
    document.getElementById('cart-total').textContent = total.toFixed(2);
}