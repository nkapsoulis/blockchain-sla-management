import classNames from 'classnames';
import React, { MouseEventHandler } from 'react';
import './Button.css';

interface Props {
  onClick?: MouseEventHandler,
  children: JSX.Element | string,
  fullWidth?: boolean,
  className?: string,
  disabled?: boolean,
}

function Button({
  onClick,
  disabled,
  fullWidth,
  children,
  className,
}:Props) {
  const classes = classNames(
    'btn',
    className,
    {
      'full-width': fullWidth,
    },
  );

  return (
    <button type={onClick ? 'button' : 'submit'} className={classes} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  );
}

Button.defaultProps = {
  className: '',
  fullWidth: false,
  disabled: false,
  onClick: undefined,
};

export default Button;
