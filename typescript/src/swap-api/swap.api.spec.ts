import {SwapApi} from './swap.api'
import {SwapRequest} from './swap.request'
import {SwapRequestParams} from './types'
import {HttpProviderConnector} from '../connectors';

describe('Swap API', () => {
    const httpProvider: HttpProviderConnector = {
        get: jest.fn().mockImplementationOnce(() => {
            return Promise.resolve()
        }),
        post: jest.fn().mockImplementation(() => {
            return Promise.resolve()
        })
    }

    it('should submit one order', async () => {
        const quoter = SwapApi.new(
            {
                url: 'https://test.com/relayer',
                network: 1
            },
            httpProvider
        )

        const orderData: SwapRequestParams = {
            // order: {
            //     allowedSender: '0x0000000000000000000000000000000000000000',
            //     interactions:
            //         '0x000c004e200000000000000000219ab540356cbb839cbe05303d7705faf486570009',
            //     maker: '0x00000000219ab540356cbb839cbe05303d7705fa',
            //     makerAsset: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
            //     makingAmount: '1000000000000000000',
            //     offsets: '0',
            //     receiver: '0x0000000000000000000000000000000000000000',
            //     salt: '45118768841948961586167738353692277076075522015101619148498725069326976558864',
            //     takerAsset: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
            //     takingAmount: '1420000000'
            // },
            signature: '0x123signature-here789',
            quoteId: '9a43c86d-f3d7-45b9-8cb6-803d2bdfa08b'
        }

        const params = SwapRequest.new(orderData)

        await quoter.submit(params)

        expect(httpProvider.post).toHaveBeenCalledWith(
            'https://test.com/relayer/v1.0/1/order/submit',
            orderData
        )
    })

    it('should submit two orders order', async () => {
        const quoter = SwapApi.new(
            {
                url: 'https://test.com/relayer',
                network: 1
            },
            httpProvider
        )

        const orderData1: SwapRequestParams = {
            // order: {
            //     allowedSender: '0x0000000000000000000000000000000000000000',
            //     interactions:
            //         '0x000c004e200000000000000000219ab540356cbb839cbe05303d7705faf486570009',
            //     maker: '0x00000000219ab540356cbb839cbe05303d7705fa',
            //     makerAsset: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
            //     makingAmount: '1000000000000000000',
            //     offsets: '0',
            //     receiver: '0x0000000000000000000000000000000000000000',
            //     salt: '45118768841948961586167738353692277076075522015101619148498725069326976558864',
            //     takerAsset: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
            //     takingAmount: '1420000000'
            // },
            signature: '0x123signature-here789',
            quoteId: '9a43c86d-f3d7-45b9-8cb6-803d2bdfa08b'
        }

        const orderData2: SwapRequestParams = {
            // order: {
            //     allowedSender: '0x0000000000000000000000000000000000000000',
            //     interactions:
            //         '0x000c004e200000000000000000219ab540356cbb839cbe05303d7705faf486570009',
            //     maker: '0x12345678219ab540356cbb839cbe05303d771111',
            //     makerAsset: '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
            //     makingAmount: '1000000000000000000',
            //     offsets: '0',
            //     receiver: '0x0000000000000000000000000000000000000000',
            //     salt: '45118768841948961586167738353692277076075522015101619148498725069326976558864',
            //     takerAsset: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
            //     takingAmount: '1420000000'
            // },
            signature: '0x123signature-2-here789',
            quoteId: '1a36c861-ffd7-45b9-1cb6-403d3bdfa084'
        }

        const params = [
            SwapRequest.new(orderData1),
            SwapRequest.new(orderData2)
        ]

        await quoter.submitBatch(params)

        expect(httpProvider.post).toHaveBeenCalledWith(
            'https://test.com/relayer/v1.0/1/order/submit/many',
            params
        )
    })
})
