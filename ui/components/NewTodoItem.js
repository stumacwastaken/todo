import { Card, CardBody, Flex, Icon, IconButton, Input, InputGroup, InputRightElement } from "@chakra-ui/react";
import { useState } from "react";
import { MdCancel } from 'react-icons/md';
export default function NewTodoItem({onCreate}) {
    const [value, setValue] = useState("")
    const handleChange = (event) => {
        setValue(event.target.value)
    }
    const handleEnter = (event) => {
        if (event.key === "Enter"){
           onCreate(value)
           setValue("")
        }
    }
    const inputter = <Input value={value} onChange={handleChange} onKeyUp={handleEnter} variant="unstyled" placeholder="New Todo..." size="lg"  _placeholder={{ opacity: 1, color: 'gray.200' }}/>

  
    const onClear = (event) => {
        setValue("")
    }
    return (
        <Card mb={3} height="80px">
            <CardBody height="100%">
                <Flex height="100%" width="100%" justifyContent="space-between">
                <InputGroup>
               {inputter}
                <InputRightElement children={ <IconButton variant="ghost" aria-label='Search database' icon={<Icon as={MdCancel} />} onClick={onClear} />}/>          \
                </InputGroup>
                </Flex>
            </CardBody>

        </Card>

    )
}

/**<InputRightElement children={<CheckIcon color='green.500' />} /> */