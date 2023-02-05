import { Card, CardBody, Flex, Text } from '@chakra-ui/react'

export default function FailedCard({ error }) {

    return (
        <Card mb={3} maxH="80px">
            <CardBody height="100%" bgColor="red.200">
                <Flex width="100%" justifyContent="space-between">
                    <Text alignSelf="center"
                        textDecorationStyle="solid" overflow="clip">
                        <span> {error.message} </span>
                    </Text>

                </Flex>
            </CardBody>
        </Card>
    )
}