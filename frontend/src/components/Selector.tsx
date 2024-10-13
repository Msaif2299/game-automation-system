import React, { useEffect, useState } from "react";
import Select, { ControlProps, CSSObjectWithLabel, GroupBase, StylesConfig } from "react-select";
import { capitalize, useDebounce } from "../common/common";

export function SelectorLabel({labelText}: {labelText: string}) {
    return <label className="glitch" data-text={labelText}>{labelText}</label>
}

interface SelectorProps {
    name: string
    scripts: string[]
    loadOptions: () => Promise<void>
    handleOnChange: (selectedOptions: string) => void
}

interface OptionType {
    label: string
    value: string
}

export default function Selector({name, scripts, loadOptions, handleOnChange} : SelectorProps) {
    const [options, setOptions] = useState<OptionType[]>([])
    const debouncedOnChange = useDebounce((selectedOption: OptionType) => {
        handleOnChange(selectedOption.value)
    }, 300)

    const customStyles: StylesConfig<OptionType, false> = {
        control: (
            provided: CSSObjectWithLabel, 
            state: ControlProps<OptionType, false, GroupBase<OptionType>>
        ) => ({
            ...provided, // Spread the provided styles first
            backgroundColor: state.isFocused ? '#000000' : '#000000', // Change background color based on focus state
            borderColor: state.isFocused ? '#3abd40' : '#3abd40', // Change border color based on focus state
            boxShadow: state.isFocused ? '0 0 0 1px #3abd40' : undefined, // Apply box-shadow when focused
            color: '#3abd40',
            '&:hover': {
                borderColor: state.isFocused ? '#3abd40' : '#000000', // Adjust hover behavior
            },
        }),
        menu: (provided) => ({
            ...provided,
            backgroundColor: '#000000', // Change dropdown menu background color
            borderColor: '#3abd40',
            borderWidth: '1px',
            borderStyle: 'solid',
        }),
        menuList: (provided) => ({
            ...provided,
            borderColor: '#3abd40'
        }),
        option: (provided, state) => ({
            ...provided,
            backgroundColor: state.isSelected ? '#000000' : state.isFocused ? '#000000' : '#000000',
            color: state.isSelected ? '#3abd40' : '#3abd40',
            ':active': {
                ...provided[':active'],
                backgroundColor: state.isSelected ? '#000000' : '#000000',
            },
            borderBottom: '1px solid #3abd40',
            ':last-of-type': {
                borderBottom: 'none',
            }
        }),
        singleValue: (provided) => ({
            ...provided,
            color: '#3abd40'
        }),
        placeholder: (provided) => ({
            ...provided,
            color: '#3abd40'
        }),
        dropdownIndicator: (provided) => ({
            ...provided,
            color: '#3abd40'
        }),
        indicatorSeparator: (provided) => ({
            ...provided,
            backgroundColor: '#3abd40',
        }),
    };

    useEffect(() => {
        loadOptions().then(() => {
            setOptions(scripts.map((script) => ({ label: capitalize(script), value: script})))
        });
    }, [loadOptions, scripts]);
    return (
        <Select 
            className="selectors"
            name={name}
            options={options}
            onChange={(selectedOption) => debouncedOnChange(selectedOption as OptionType)}
            styles={customStyles}
        />
    );
}