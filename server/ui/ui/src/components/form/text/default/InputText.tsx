// vendor
import React, { FC, KeyboardEvent, RefObject } from 'react';
import classNames from 'classnames';
// css
import './InputText.scss';


interface Props {
  css?: string;
  disabled?: boolean;
  inputRef: RefObject<HTMLInputElement>;
  isRequired?: boolean;
  flexParent?: boolean;
  label: string;
  placeholder?: string;
  isPassword?: boolean;
  updateValue: any;
  onFocus?: any;
  onBlur?: any;
  onChange?: any;
  // TODO: debug why this fails with search implementation
  // updateValue: (event: KeyboardEvent<HTMLInputElement>) => void;
}

const InputText: FC<Props> = ({
  css = null,
  disabled = false,
  inputRef,
  isRequired,
  flexParent,
  label,
  placeholder,
  updateValue,
  isPassword,
  onFocus = null,
  onBlur = null,
  onChange = null,
}:Props) => {

  // declare css here
  const inputTextCSS = classNames({
    InputText__input: true,
    [css]: typeof css === 'string'
  });
  const inputParentCSS = classNames({
    InputText: true,
    'flex-1': flexParent,
  });
  return (
    <div className={inputParentCSS}>
      <label className="InputText__label">
        {label}
        { isRequired &&
          <i className="InputText__i InputText__i--orange">* Required</i>
        }
      </label>
      <input
        className={inputTextCSS}
        disabled={disabled}
        onKeyUp={(event) => {updateValue(event)}}
        placeholder={placeholder || ''}
        ref={inputRef}
        type={isPassword ? 'password': 'text'}
        onFocus={onFocus}
        onBlur={onBlur}
        onChange={onChange}
      />
    </div>
  )
}


export default InputText;
