import React from 'react';

import names from '../../utils/names';

import './Input.css';

interface Props {
  type: string,
  id: string,
  name: string,
  placeholder?: string,
  required?: boolean,
  value?: string,
  onChange?: React.ChangeEventHandler<HTMLInputElement> |
  React.ChangeEventHandler<HTMLTextAreaElement>,

}

function Input({
  type, id, name, required, placeholder, onChange, value,
}: Props) {
  const classes = 'input';
  return (
    <label htmlFor={id}>
      {name}
      { type === 'textarea'
        ? (
          <textarea
            id={id}
            name={names.kebabCase(name)}
            required={required}
            className={classes}
            placeholder={placeholder}
            onChange={onChange as React.ChangeEventHandler<HTMLTextAreaElement>}
            value={value}
          />
        )

        : (
          <input
            type={type}
            id={id}
            required={required}
            name={names.kebabCase(name)}
            className={classes}
            placeholder={placeholder}
            onChange={onChange as React.ChangeEventHandler<HTMLInputElement>}
            value={value}
          />
        )}
    </label>
  );
}

Input.defaultProps = {
  required: false,
  placeholder: '',
  value: undefined,
  onChange: undefined,
};

export default Input;
