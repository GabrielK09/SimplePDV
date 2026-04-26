interface CustomerContract {
    readonly id: number;
    name: string;
    cpf_cnpj: string;
    deleted_at?: Date;
};
