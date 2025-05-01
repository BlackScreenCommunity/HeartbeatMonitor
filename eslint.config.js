import eslintPluginTs from '@typescript-eslint/eslint-plugin';
import eslintParserTs from '@typescript-eslint/parser';

/** @type {import('eslint').Linter.Config} */
export default [{
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
        parser: eslintParserTs,
        parserOptions: {
            project: './tsconfig.json',
            ecmaVersion: 2020,
            sourceType: 'module',
        },
    },
    plugins: {
        '@typescript-eslint': eslintPluginTs,
    },
    rules: {
        ...eslintPluginTs.configs.recommended.rules,
        "no-debugger": "error"
    },
}];
