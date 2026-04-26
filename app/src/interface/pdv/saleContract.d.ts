type StatusSale = 'Pendente'|'Concluída'|'Cancelada';

interface SaleContract {
    id: number|null;
    customer_id: number;
    sale_value?: number;
    customer: string;
    specie: string;
    products: SaleItemContract[];
    status?: StatusSale;
};

interface SaleItemContract {
    id: number|null;
    product_id: number|null;
    name: string|null;
    sale_value: number|null;
    qtde: number;
    use_grid?: boolean;
    product_with_characteristics: ProductCharacteristicsContract[]
};

interface PaySaleContract {
    sale_id: number;
    specie: 'Dinheiro'|'Pix';
    amount_paid: number
};

interface PDVContract {
    readonly id: number;
    customer: string;
    specie: string;
    sale_value: number;
    status: StatusSale;
};
