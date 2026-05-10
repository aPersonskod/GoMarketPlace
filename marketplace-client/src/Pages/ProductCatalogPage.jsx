import '../App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import Products from "../Products.jsx";
import ProductCart from "../ProductCart.jsx";
import Button from "react-bootstrap/Button";
import { useNavigate } from 'react-router';
import {useEffect, useState} from "react";
import {ApiHelper} from "../ApiHelper.jsx";

const ProductCatalogPage = () => {
    const navigate = useNavigate();
    const [cart, setCart] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [ammountToPay, setAmmountToPay] = useState(0);
    let apiHelper = new ApiHelper();

    const fetchCartData = async () => {
        try {
            let token = apiHelper.getAccessToken();
            if (token === null) return;
            let query = `${apiHelper.orderServiceBaseAddress}/get-cart`;
            let options = {
                method: 'GET',
                headers: {
                    'Authorization' : `Bearer ${token}`
                }
            };
            const response = await fetch(query, options);
            if (!response.ok) {
                setCart(null);
                return;
                throw new Error(`HTTP error! status: ${response.status}`);
            }   
            const cartResult = await response.json();
            //if(cartResult.isConfirmed) setOrders([]);
            setCart(cartResult);
        } catch (err) {
            setError(err);
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteCart = async () => {
        try {
            let token = apiHelper.getAccessToken();
            let query = `${apiHelper.orderServiceBaseAddress}/delete-cart/${cart.id}`;
            const response = await fetch(query, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization' : `Bearer ${token}`
                }//,
                //body: JSON.stringify(requestBody),
            });

            if (!response.ok) {
                //throw new Error(`HTTP error! status: ${response.status}`);
                alert(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            console.log('Update successful:', data);
            await fetchCartData();
        } catch (err) {
            setError(err);
        } finally {
            setLoading(false);
        }
    }
    
    useEffect(() => {
        fetchCartData();
    }, []);
    const handleConfirmation = () => {
        if(cart !== null){
            navigate('confirmation');
        }
    }

    if (loading) return <div>Loading data...</div>;
    if (error) return <div>Error: {error.message}</div>;
    
    return (
        <>
            <div className='row gy-2'>
                <div className='col col-xs-12 col-sm-12 col-lg-6'>
                    <p className='fs24'>Продукты:</p>
                    <div className='divStyle xsDivStyle mdDivStyle'>
                        <Products cart={cart} setAmmountToPay={setAmmountToPay} refreshCartFunc={fetchCartData}/>
                    </div>
                </div>
                <div className='col col-xs-12 col-sm-12 col-lg-6'>
                    <div className='d-flex'>
                        <p className='fs24'>Корзина:</p>
                        {cart !== null &&
                            <div>
                                <Button style={{backgroundColor:"red", marginLeft:"10px"}} onClick={handleDeleteCart}>Удалить</Button>
                            </div>
                        }
                        {cart !== null &&
                            <div className='d-flex'>
                                <p className='fs24' style={{marginLeft:"10px"}}>Сумма к оплате:</p>
                                <p className='fs24' style={{marginLeft:"10px"}}>{cart.amountToPay}</p>
                            </div>
                        }   
                    </div>
                    <div className='divStyle xsDivStyle mdDivStyle'>
                        <ProductCart key={ammountToPay} cart={cart}/>
                    </div>
                </div>
                {((cart !== null) && (cart.ammountToPay !== 0)) &&
                <div className='col col-xs-12 col-sm-12'>
                    <Button onClick={handleConfirmation}>Подтвердить заказ</Button>
                </div>
                }
            </div>
        </>
    );
}
export default ProductCatalogPage;