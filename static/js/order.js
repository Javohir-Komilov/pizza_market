function viewOrderDetails(orderId) {
    // Здесь вы можете реализовать логику для отображения деталей заказа
    // Например, открыть модальное окно или перенаправить на страницу с деталями заказа
    console.log("Viewing details for order", orderId);
}

function cancelOrder(orderId) {
    if (confirm("Are you sure you want to cancel this order?")) {
        fetch(`/user/orders/${orderId}/cancel`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert("Order cancelled successfully");
                location.reload(); // Перезагрузить страницу для обновления статуса
            } else {
                alert("Failed to cancel order: " + data.message);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert("An error occurred while cancelling the order");
        });
    }
}

function reorder(orderId) {
    fetch(`/user/orders/${orderId}/reorder`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert("Order placed successfully");
            window.location.href = "/user/cart"; // Перенаправить в корзину
        } else {
            alert("Failed to place order: " + data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert("An error occurred while placing the order");
    });
}