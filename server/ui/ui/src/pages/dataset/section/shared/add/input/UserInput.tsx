// vendor
import React,
{
  FC,
  KeyboardEvent,
  RefObject,
} from 'react';
// components
import InputText from 'Components/form/text/index';
// css
import './UserInput.scss';


interface Props {
  inputRef: RefObject<HTMLInputElement>;
  permissionType: string;
  updateName: (event: KeyboardEvent) => void;
  checkName: any;
}

const AddPermission: FC<Props> = ({
  inputRef,
  permissionType,
  updateName,
  checkName,
}:Props) => {


  return (
    <div className="UserInput">
      <InputText
        inputRef={inputRef}
        label={permissionType === 'user' ? 'Add a user by username or email address' : `Add ${permissionType}`}
        placeholder={permissionType === 'user' ? 'username or email': `${permissionType}name`}
        updateValue={updateName}
        onChange={checkName}
      />
    </div>
  );
}


export default AddPermission;
