# Kyora Prompt Files

Reusable prompt files for common development tasks in the Kyora monorepo.

## Available Prompts

### Feature Development

- **`/add-backend-feature`** - Add new feature to backend (Go API)
- **`/improve-backend-feature`** - Improve an existing feature in the backend (Go API)
- **`/add-portal-feature`** - Add new feature to portal-web (React dashboard)
- **`/improve-portal-feature`** - Improve an existing feature in portal-web (React dashboard)

### Cross-Project

- **`/cross-project-enhancement`** - Enhancement spanning backend + portal-web
- **`/cross-project-refactor`** - Refactor spanning backend + portal-web

### Project Scaffolding

- **`/create-backend-project`** - Create new backend domain module
- **`/create-frontend-project`** - Create new frontend feature module

### Debugging

- **`/debug-backend`** - Debug and fix backend issues
- **`/debug-portal-web`** - Debug and fix portal-web issues
- **`/fix-ui-ux`** - Fix UI/UX issues in portal-web

### Refactoring

- **`/refactor-backend-feature`** - Refactor a single backend feature/domain
- **`/refactor-backend-project`** - Project-wide refactor in the backend
- **`/refactor-portal-feature`** - Refactor a single portal-web feature
- **`/refactor-portal-project`** - Project-wide refactor in portal-web

## Usage

### In Chat View

Type `/` followed by the prompt name in the chat input:

```
/add-backend-feature Add endpoint to export orders as CSV
```

```
/fix-ui-ux Button alignment broken in RTL mode on mobile
```

### From Command Palette

1. Open Command Palette (`⌘⇧P` / `Ctrl+Shift+P`)
2. Run `Chat: Run Prompt`
3. Select prompt from list
4. Enter required parameters

### From Prompt File

1. Open prompt file in editor
2. Click play button in title bar
3. Choose to run in current or new chat session

## Prompt Structure

Each prompt file includes:

- **YAML Frontmatter**: Description, agent, tools, model configuration
- **Feature Requirements**: Input variable for user to describe their task
- **Instructions**: Links to relevant instruction files
- **Standards**: Quality and architecture requirements
- **Workflow**: Step-by-step implementation guide
- **Done Criteria**: Completion checklist

## Tips

1. **Be Specific**: Provide detailed feature/issue descriptions in the input
2. **Reference Examples**: Prompts suggest similar existing features to reference
3. **Follow Standards**: Prompts enforce Kyora architecture patterns
4. **Test Thoroughly**: Prompts include testing steps in workflow
5. **Use Instruction Files**: Prompts link to specialized instruction files for details

## Instruction Files Reference

Prompts reference these instruction files (SSOT):

- `backend-core.instructions.md` - Backend architecture
- `backend-testing.instructions.md` - Backend testing
- `portal-web-architecture.instructions.md` - Frontend architecture
- `portal-web-development.instructions.md` - Frontend workflow
- `forms.instructions.md` - Form system
- `ui-implementation.instructions.md` - UI patterns
- `design-tokens.instructions.md` - Design tokens
- `charts.instructions.md` - Data visualization
- `ky.instructions.md` - HTTP client
- `stripe.instructions.md` - Billing
- `resend.instructions.md` - Email
- `asset_upload.instructions.md` - File uploads

## Adding New Prompts

1. Create `.prompt.md` file in this directory
2. Add YAML frontmatter with description, agent, tools
3. Use `${input:variableName:placeholder}` for user inputs
4. Reference instruction files with relative links
5. Include workflow steps and done criteria
6. Update this README with new prompt
7. Test prompt with various inputs

## Examples

Below are **small, complete** examples for every prompt.

Tip: you can run a prompt by typing `/prompt-name` in chat, then filling in the requested inputs. If a prompt has **Constraints**, they are optional and should be written as hard rules.

### Feature Development

#### `/add-backend-feature`

When asked for **Feature Requirements**:

```
Add an endpoint to export orders as CSV for a business.
Include filters: date range + status.
```

#### `/improve-backend-feature`

When asked:

```
Improvement Brief:
Reduce N+1 queries on the orders list endpoint.

Constraints (optional):
No DB migrations. No endpoint changes. No response shape changes.
```

#### `/add-portal-feature`

When asked for **Feature Requirements**:

```
Add order filtering by status with date range to the orders list.
Must work for both Arabic and English.
```

#### `/improve-portal-feature`

When asked:

```
Improvement Brief:
Improve inventory adjustment UX on mobile: clearer errors + better loading states.

Constraints (optional):
No new routes. No backend changes. Keep existing layout.
```

### Cross-Project

#### `/cross-project-enhancement`

When asked:

```
Enhancement Brief:
Add a new computed field to the order API response (profit) and display it in the order details page.

Constraints (optional):
No DB migrations. Keep existing endpoints; only extend the response.
```

#### `/cross-project-refactor`

When asked:

```
Refactor Brief:
Unify backend problem error codes for order create/update and update portal form error mapping accordingly.

Constraints (optional):
No behavior changes. Portal UI should look the same.
```

### Project Scaffolding

#### `/create-backend-project`

When asked for **Project Requirements**:

```
subscription
```

#### `/create-frontend-project`

When asked for **Project Requirements**:

```
subscriptions management
```

### Debugging

#### `/debug-backend`

When asked for **Issue Description**:

```
Getting 500 when creating an order with multiple line items; error occurs after inventory is updated.
```

#### `/debug-portal-web`

When asked for **Issue Description**:

```
Order create form submits, but the UI stays loading forever and no error is shown.
```

#### `/fix-ui-ux`

When asked for **Issue Description**:

```
Button alignment is broken in RTL mode on mobile in the orders list toolbar.
```

### Refactoring

#### `/refactor-portal-feature`

When asked:

```
Refactor Brief:
Extract reusable OrderForm fields into shared components and unify Zod schema reuse.

Feature/Area:
orders

Constraints (optional):
No UI changes. No backend changes.
```

#### `/refactor-portal-project`

When asked:

```
Refactor Brief:
Convert POST/PUT/DELETE requests to TanStack Query mutations across the app.

Constraints (optional):
No UI changes. No route changes. No new dependencies.
```

#### `/refactor-backend-feature`

When asked:

```
Refactor Brief:
Standardize RFC 7807 problems for inventory adjustments and remove duplicated validation code.

Domain/Area:
inventory

Constraints (optional):
No endpoint changes. No DB migrations. No behavior changes.
```

#### `/refactor-backend-project`

When asked:

```
Refactor Brief:
Standardize request validation + problem errors across handlers.

Constraints (optional):
No endpoint changes. No DB migrations. Must be mechanical and low-risk.
```
