interface CashRegisterContract {
    readonly id: number;
    sale_id: number;
    shopping_id: number;
    description: string;
    customer_id: unknown;
    customer: string;
    specie_id: number;
    specie: string;
    input_value: number;
    output_value: number;
    total_balance: number;
};
