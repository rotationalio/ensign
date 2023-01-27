import { vi } from 'vitest';
import axios from "axios";
import { createAccountRequest } from "../AccountApiService";
import axiosInstance from '@/application/api/ApiService';

describe("AccountApiService", () => {
    describe("createAccountRequest", () => {
        it("returns request resolved with response", async () => {
            const mockAccount = { email: "test@ensign.com", password: "test", username: "test" }

            const requestSpy = vi.fn().mockReturnValueOnce({
                status: 200,
                data: mockAccount,
                statusText: "OK"
            });
            const request = createAccountRequest(requestSpy);
            const response = await request(mockAccount);
            expect(response).toBe(mockAccount);
            expect(requestSpy).toHaveBeenCalledTimes(1);

            // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/register', {
            //     method: "POST",
            //     headers: {
            //         "Content-Type": "application/json",
            //     },
            //     body: JSON.stringify(mockAccount),
            // });

            // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/register');
            // expect(requestSpy).toHaveBeenCalledWith(expect.stringContaining('localhost:8088'));

        });
    });
});