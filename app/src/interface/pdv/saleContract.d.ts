interface SaleContract {
    readonly id: number;
    readonly customer_id: number;
    customer: string;
    specie: string;
    products: SaleItemContract[]
};

interface SaleItemContract {
    readonly id: number;
    product_id: number;
    name: string;
    price: number;
    qtde: number;
};

interface PaySaleContract {
    sale_id: number;
    specie: 'Dinheiro'|'Pix';
    amount_paid: number
};

type StatusSale = 'Pendente'|'Concluída'|'Cancelada';

interface PDVContract {
    readonly id: number;
    customer: string;
    specie: string;
    sale_value: number;
    status: StatusSale;
};
