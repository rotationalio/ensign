import { describe, expect, it, vi } from "vitest";
import { createProjectTopic} from '../createTopicApiService'

describe('Topic', () => {
    describe('Create Topic', () => {
        it('returns request with a response', async () => {
            const mockTopic = {
                id: '1',
                name: 'topic01'
            }
            const requestSpy = vi.fn().mockReturnValueOnce({
                status: 200,
                data: mockTopic,
                statusText: 'OK',
            });
            const mockDTO = {
                name: 'topic01'
            };
            const request = createProjectTopic(requestSpy)
            const response = await request(mockDTO)
            expect(response).toBe(mockTopic)
            expect(requestSpy).toHaveBeenCalledTimes(1)
        });
        it('returns request with error', async () => {
            const requestSpy = vi.fn().mockReturnValueOnce({
                status: 200,
                data: undefined,
                statusText: 'OK',
            });
            const mockDTO = {
                name: 'topic01'
            };
            const request = createProjectTopic(requestSpy)
            const response = await request(mockDTO)
            expect(response).toBe(undefined)
            expect(requestSpy).toHaveBeenCalledTimes(1)
        });
    });
})