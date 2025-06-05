/* eslint-disable @typescript-eslint/unbound-method */
/* eslint-disable @typescript-eslint/no-unsafe-return */
import { Test, TestingModule } from '@nestjs/testing';
import { Request } from 'express';
import { LineStrategy } from './line.strategy';
import { VerifyIDToken } from '@line/bot-sdk';
import { UnauthorizedException } from '@nestjs/common';
import { LineService } from 'src/core/line/line.service';

describe('LineStrategy', () => {
  let lineStrategy: LineStrategy;
  let lineService: LineService;
  let mockRequest: Partial<Request>;
  let mockSuccess: jest.Mock;
  let mockFail: jest.Mock;

  const mockUser: VerifyIDToken = {
    sub: 'U1234567890',
    name: 'Test User',
    picture: 'https://example.com/picture.jpg',
    email: 'test@example.com',
    scope: '',
    client_id: '',
    expires_in: 0,
    iss: '',
    aud: 0,
    exp: 0,
    iat: 0,
    nonce: '',
    amr: [],
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        LineStrategy,
        {
          provide: LineService,
          useValue: {
            verifyIdToken: jest.fn(),
          },
        },
      ],
    }).compile();

    lineStrategy = module.get<LineStrategy>(LineStrategy);
    lineService = module.get<LineService>(LineService);

    // Mock PassportStrategy methods
    mockSuccess = jest.fn();
    mockFail = jest.fn();
    lineStrategy['success'] = mockSuccess;
    lineStrategy['fail'] = mockFail;

    // Mock Request
    mockRequest = {
      headers: {
        authorization: 'Bearer mock-id-token',
      },
    };
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('authenticate', () => {
    it('should call success when ID token is valid', async () => {
      // Arrange
      jest
        .spyOn(lineService, 'verifyIdToken')
        .mockReturnValue(Promise.resolve(mockUser));

      // Act
      lineStrategy.authenticate(mockRequest as Request);

      // Assert
      await new Promise(process.nextTick); // รอ async operation
      expect(lineService.verifyIdToken).toHaveBeenCalledWith('mock-id-token');
      expect(mockSuccess).toHaveBeenCalledWith(mockUser);
      expect(mockFail).not.toHaveBeenCalled();
    });

    it('should call fail when no ID token is provided', async () => {
      // Arrange
      mockRequest.headers = {};

      // Act
      lineStrategy.authenticate(mockRequest as Request);

      // Assert
      await new Promise(process.nextTick);
      expect(lineService.verifyIdToken).not.toHaveBeenCalled();
      expect(mockFail).toHaveBeenCalledWith('No ID token provided', 401);
      expect(mockSuccess).not.toHaveBeenCalled();
    });

    it('should call fail when ID token is invalid', async () => {
      // Arrange
      jest
        .spyOn(lineService, 'verifyIdToken')
        .mockReturnValue(
          Promise.reject(new UnauthorizedException('Invalid LINE ID Token')),
        );

      // Act
      lineStrategy.authenticate(mockRequest as Request);

      // Assert
      await new Promise(process.nextTick);
      expect(lineService.verifyIdToken).toHaveBeenCalledWith('mock-id-token');
      expect(mockFail).toHaveBeenCalledWith('Invalid LINE ID Token', 401);
      expect(mockSuccess).not.toHaveBeenCalled();
    });
  });

  describe('validate', () => {
    it('should return payload with correct structure', () => {
      // Act
      const result = lineStrategy.validate(mockUser);

      // Assert
      expect(result).toEqual({
        sub: 'U1234567890',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email: 'test@example.com',
        scope: '',
        client_id: '',
        expires_in: 0,
        iss: '',
        aud: 0,
        exp: 0,
        iat: 0,
        nonce: '',
        amr: [],
      });
    });
  });
});
