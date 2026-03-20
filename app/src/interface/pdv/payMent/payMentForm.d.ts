/**
 * Para create
 */
interface PayMentFormContract {
    readonly id: number;
    specie: string;
    pix_key: string;
};

/**
 * Para o pagamento na api
 */
interface PayMentPayLoadContract {
    sale_id?: number;
    shopping_id?: number;
    species: any[];
};

/**
 * Para venda
 */
interface PayMentValue {
    id: number;
    specie: string;
    amount: string;
    pix_key: string;
};

/**
 * Usado para cancelamento tanto para venda quanto para compra
 */
interface CancelContract {
    shopping_id?: number;
    sale_id?: number;
    route: string;
};