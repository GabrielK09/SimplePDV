interface SaleContract {
    id: number|null;
    customer_id: number;
    customer: string;
    specie: string;
    products: SaleItemContract[]
};

interface SaleItemContract {
    id: number|null;
    product_id: number|null;
    name: string;
    price: number|null;
    qtde: number;
    product_with_characteristics: ProductCharacteristicsContract|null
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
