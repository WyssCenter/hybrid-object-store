// vendor
import React, {
    FC,
    ChangeEvent,
    useCallback,
    useRef,
  } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSquare } from '@fortawesome/free-regular-svg-icons';
import { faCheckSquare } from '@fortawesome/free-solid-svg-icons';
// css
import './Checkbox.scss';


interface Props {
  disabled?: boolean;
  id: string;
  isChecked: boolean;
  label: string;
  updateCheckbox: (list: any) => void;
}

const Checkbox: FC<Props> = ({
  disabled = false,
  id,
  isChecked,
  label,
  updateCheckbox,
}:Props) => {
  //refs
  const inputRef = useRef(null);

  // vars
  const icon = isChecked
    ? faCheckSquare
    : faSquare;

  /**
  * Method updates the search value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setSearch}
  */
  const handleChangeEvent = useCallback((event: ChangeEvent) => {
    const element = event.target as HTMLInputElement;

    updateCheckbox(element.checked)
    return;
  }, [updateCheckbox])

    const root = document.documentElement;
    const primaryHash = root.style.cssText ?
      root.style.cssText.split(';')[0].split(':')[1]
      : '#957299';

  return (
    <div className="Checkbox">
      <label
        className="Checkbox"
        htmlFor={id}
      >
        <FontAwesomeIcon
          color={primaryHash}
          icon={icon}
          size="lg"
        />
        <span>{label}</span>
      </label>
      <input
        className="hidden"
        disabled={disabled}
        id={id}
        onChange={(event) => handleChangeEvent(event)}
        ref={inputRef}
        type="checkbox"
      />
    </div>
  )
}


export default Checkbox;
