#include "cal_30.h"
#include "ui_cal_30.h"

Cal_30::Cal_30(QString uid, QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::Cal_30), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);

    ui->events_scroll->setWidgetResizable(true);
    on_calendarWidget_clicked(QDate::currentDate());
}

Cal_30::~Cal_30()
{
    delete ui;
}

void Cal_30::on_calendarWidget_clicked(const QDate &date)
{
    {
        QWidget* prev = ui->events_scroll->widget();
        if (prev != nullptr) {
            delete prev;
        }
    }
    QWidget *central = new QWidget;
    QVBoxLayout* layout = new QVBoxLayout(central);
    QVector<EventData> events = this->cli.getEventsByDay(date, uid);
    for (auto e : events) {
        event_entry *item = new event_entry(e, uid);
        connect(item, &event_entry::deleted, this, &Cal_30::eventDeleted);
        connect(item, &event_entry::edited, this, &Cal_30::eventEdited);
        layout->addWidget(item);
    }
    layout->addStretch();
    central->setSizePolicy(QSizePolicy::Preferred, QSizePolicy::MinimumExpanding);
    ui->events_scroll->setWidget(central);
}

void Cal_30::eventDeleted(event_entry* ev) {
    bool ok = cli.deleteEvent(uid, ev->getData().ID);
    if (!ok) {
        QMessageBox::warning(this, "Ошибка", "При удалении события произошла ошибка");
    } else {
        delete ev;
    }
    this->on_calendarWidget_clicked(ui->calendarWidget->selectedDate());
}
void Cal_30::eventEdited(event_entry* ev) {
    on_calendarWidget_clicked(ui->calendarWidget->selectedDate());
}
