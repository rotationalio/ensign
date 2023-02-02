export interface APIKey {
	id: string;
    client_id: string;
    client_secret: string;
    name: string;
    owner: string;
    permissions: string[];
    created: string;
    modified: string;
}

export type NewAPIKey = Omit<APIKey, 'id'>;