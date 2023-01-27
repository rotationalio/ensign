import { vi } from 'vitest';
import { loginRequest } from "../LoginApiService";

describe("LoginApiService", () => {
    describe("LoginRequest", () => {
        it("returns request resolved with response", async () => {
            const mockUser = { password: "test", username: "test" };
            const mockResponse = { access_token: 'xxx', refresh_token: 'yyy' }

            const requestSpy = vi.fn().mockReturnValueOnce({
                status: 200,
                data: mockResponse,
                statusText: "OK"
            });
            const request = loginRequest(requestSpy);
            const response = await request(mockUser);
            expect(response).toBe(mockResponse);
            expect(requestSpy).toHaveBeenCalledTimes(1);

            // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/login', {
            //     method: "POST",
            //     headers: {
            //         "Content-Type": "application/json",
            //     },
            //     body: JSON.stringify(mockUser),
            // });

            // expect(requestSpy).toHaveBeenCalledWith('http://localhost:8088/v1/login');
            // expect(requestSpy).toHaveBeenCalledWith(expect.stringContaining('localhost:8088'));

        });
    });
});