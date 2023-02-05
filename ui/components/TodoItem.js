import { Card, CardBody, Checkbox, Flex, Icon, IconButton, Text } from "@chakra-ui/react";
import { MdArchive } from 'react-icons/md';
export default function TodoItem({ item, onComplete, onArchive }) {
    return (
        <Card mb={3} maxH="80px">
            <CardBody height="100%">
                {/* //TODO: get rid of space-between and make it more like keep */}
                <Flex width="100%" justifyContent="space-between">
                    <Checkbox isChecked={item.completed} onChange={() => onComplete(item.id, !item.completed)}></Checkbox>
                    <Text alignSelf="center" textAlign="center" textDecorationLine={item.completed ? "line-through" : "none"}
                        textDecorationStyle="solid" overflow="clip">
                        <span style={{ whiteSpace: "nowrap", overflow: "clip", textOverflow: "ellipsis" }}>{item.summary}</span>
                    </Text>
                    <IconButton variant="ghost" aria-label='Search database' icon={<Icon as={MdArchive} />} onClick={() => onArchive(item.id)} />
                </Flex>
            </CardBody>
        </Card>
    )
}