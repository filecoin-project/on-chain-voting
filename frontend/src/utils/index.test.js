import { describe, expect, it } from 'vitest';
import { convertBytes } from './index';

describe('calcFromSeconds', () => {
    it('should return 1.00 KiB from 1024', () => {
        expect(convertBytes(1024)).toBe('1.00 KiB');
    });

    it('should return 1.00 MiB from 1024 * 1024', () => {
        expect(convertBytes(1024 * 1024)).toBe('1.00 MiB');
    });

    it('should return 1.00 GiB from 1024 * 1024 * 1024', () => {
        expect(convertBytes(1024 * 1024 * 1024)).toBe('1.00 GiB');
    });
});