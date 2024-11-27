const orderWidget = document.getElementById('order-widget');
const toggleButton = document.getElementById('close-button');
const openButton = document.getElementById('open-button');

if (localStorage.getItem('orderWidgetHidden') === 'true') {
    orderWidget.style.display = 'none';
    toggleButton.textContent = 'chevron_right';
    openButton.style.display = 'flex';
}

toggleButton.addEventListener('click', function() {
    const isHidden = orderWidget.style.display === 'none';

    orderWidget.style.display = isHidden ? 'block' : 'none';

    if (isHidden) {
        toggleButton.textContent = 'chevron_left';
        localStorage.setItem('orderWidgetHidden', 'false');
        openButton.style.display = 'none';
    } else {
        toggleButton.textContent = 'chevron_right';
        localStorage.setItem('orderWidgetHidden', 'true');
        openButton.style.display = 'flex'; 
    }
});

openButton.addEventListener('click', function() {
    orderWidget.style.display = 'block';
    toggleButton.textContent = 'chevron_left';
    localStorage.setItem('orderWidgetHidden', 'false');
    openButton.style.display = 'none';
});

 
function displayOrderDetails(order) {
    const orderDetailsDiv = document.getElementById('details');

    let htmlContent = `
    <h4>Order Details:</h4>
    <p><strong>Order UID:</strong> ${order.order_uid}</p>
    <p><strong>Track Number:</strong> ${order.track_number}</p>
    <p><strong>Entry:</strong> ${order.entry}</p>
    <p><strong>Locale:</strong> ${order.locale}</p>
    <p><strong>Customer ID:</strong> ${order.customer_id}</p>
    <p><strong>Internal Signature:</strong> ${order.internal_signature || 'N/A'}</p>
    <p><strong>Delivery Service:</strong> ${order.delivery_service}</p>
    <p><strong>Shard Key:</strong> ${order.shardkey}</p>
    <p><strong>SM ID:</strong> ${order.sm_id}</p>
    <p><strong>Date Created:</strong> ${new Date(order.date_created).toLocaleString()}</p>
    <p><strong>OOF Shard:</strong> ${order.oof_shard}</p>
    
    <h4>Delivery Information</h4>
    <p><strong>Name:</strong> ${order.delivery.name}</p>
    <p><strong>Phone:</strong> ${order.delivery.phone}</p>
    <p><strong>Address:</strong> ${order.delivery.address}, ${order.delivery.city}, ${order.delivery.region}, ${order.delivery.zip}</p>
    <p><strong>Email:</strong> ${order.delivery.email}</p>
    
    <h4>Payment Information</h4>
    <p><strong>Transaction:</strong> ${order.payment.transaction}</p>
    <p><strong>Request ID:</strong> ${order.payment.request_id || 'N/A'}</p>
    <p><strong>Currency:</strong> ${order.payment.currency}</p>
    <p><strong>Provider:</strong> ${order.payment.provider}</p>
    <p><strong>Amount:</strong> $${(order.payment.amount / 100).toFixed(2)}</p>
    <p><strong>Payment Date:</strong> ${new Date(order.payment.payment_dt * 1000).toLocaleString()}</p>
    <p><strong>Bank:</strong> ${order.payment.bank}</p>
    <p><strong>Delivery Cost:</strong> $${(order.payment.delivery_cost / 100).toFixed(2)}</p>
    <p><strong>Goods Total:</strong> $${(order.payment.goods_total / 100).toFixed(2)}</p>
    <p><strong>Custom Fee:</strong> $${(order.payment.custom_fee / 100).toFixed(2)}</p>
    
    <h4>Items</h4>
    `;

    order.items.forEach(item => {
        htmlContent += `
            <div>
                <p><strong>Product Name:</strong> ${item.name}</p>
                <p><strong>Brand:</strong> ${item.brand}</p>
                <p><strong>Price:</strong> $${(item.price / 100).toFixed(2)}</p>
                <p><strong>Quantity:</strong> ${item.sale}</p>
                <p><strong>Total Price:</strong> $${(item.total_price / 100).toFixed(2)}</p>
                <p><strong>Chart ID:</strong> ${item.chrt_id}</p>
                <p><strong>RID:</strong> ${item.rid}</p>
                <p><strong>NM ID:</strong> ${item.nm_id}</p>
                <p><strong>Status:</strong> ${item.status}</p>
            </div>
        `;
    });


    orderDetailsDiv.innerHTML = htmlContent;
}

const getOrderRadio = document.getElementById('get-order');
const spamOrderRadio = document.getElementById('spam-order');
const detailsDiv = document.getElementById('details');

function displaySpamOrderIds(ordersIds) {
    let htmlContent = '<h4>Spammed Orders</h4>';

    ordersIds.forEach(id => {
        htmlContent += `<p>${id}</p>`;
    });

    detailsDiv.innerHTML = htmlContent;
}


const baseApiUrl = process.env.SERVICE_BASE_URL || 'http://localhost:8080/api/v1';
const getOrderUrl = `${baseApiUrl}/order/:id`;
const spamOrderUrl = `${baseApiUrl}/order/spam?count=`;

const orderIdInput = document.getElementById('order-id-input');
const spamCountInput = document.getElementById('spam-id-input');

let lastGetOrderDetails = null; 
let lastSpamOrderIDs = null;

getOrderRadio.addEventListener('change', () => {
    displayOrderDetails(lastGetOrderDetails);
});

spamOrderRadio.addEventListener('change', () => {
    displaySpamOrderIds(lastSpamOrderIDs);
});

document.getElementById('submit-order-id').addEventListener('click', () => {
    const orderId = orderIdInput.value;
    fetchOrderById(orderId);
    orderIdInput.value = '';
});

document.getElementById('submit-spam-count').addEventListener('click', () => {
    const count = spamCountInput.value;
    fetchSpamOrders(count);
    spamCountInput.value = '';
});

async function fetchOrderById(orderId) {
    const url = getOrderUrl.replace(':id', orderId);
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const orderData = await response.json();
        lastGetOrderDetails = orderData;
        displayOrderDetails(orderData);
    } catch (error) {
        console.error('Ошибка при получении заказа:', error);
    }
}

async function fetchSpamOrders(count) {
    const url = `${spamOrderUrl}${count}`;
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const spamData = await response.json();
        console.log(spamData);
        lastSpamOrderIDs = spamData;
        displaySpamOrderIds(spamData);
    } catch (error) {
        console.error('Ошибка при получении спамов:', error);
    }
}
