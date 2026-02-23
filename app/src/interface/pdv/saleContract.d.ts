interface SaleContract {
    readonly id: number;
    readonly id: number;
    customer: string;
    specie: string;
    products: SaleItemContract[]
};

interface SaleItemContract {
    readonly id: number;
    readonly product_id: number;
    name: string;
    price: number;
    qtde: number;
};

interface PaySaleContract {
    sale_id: number;
    specie: 'Dinheiro'|'Pix';
    amount_paid: number
};

type StatusSale = 'Pendente'|'Concluída';
interface PDVContract {
    readonly id: number;
    customer: string;
    specie: string;
    sale_value: number;
    status: StatusSale;
};
